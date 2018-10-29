package binlog

import (
	"fmt"
	"proto/mysql/pbmysql"
	"strconv"

	"vitess.io/vitess/go/vt/vttablet/tabletserver/schema"

	"golang.org/x/net/context"
	"vitess.io/vitess/go/mysql"
	binlogdatapb "vitess.io/vitess/go/vt/proto/binlogdata"
	querypb "vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/sqlparser"
)

// wrapped from binlog_streamer.go

// ParsedProto used to parse event to proto
type ParsedProto struct {
	NextPos mysql.Position
	ErrInfo error
	Events  []*pbmysql.Event
	Reload  bool
}

// ParseProtoEvents used to parse event detail from binlog
func ParseProtoEvents(ctx context.Context, events <-chan mysql.BinlogEvent, seList map[string]*schema.Engine, pbMsgCh chan *ParsedProto, startPos mysql.Position, onlyPrimaryID bool) {
	var format mysql.BinlogFormat
	var gtid mysql.GTID
	var pos mysql.Position
	var autocommit = true
	// reload is used to reload mysql data
	var reload bool
	var err error
	// pbEvents for proto defined events
	var pbEvents []*pbmysql.Event

	// Remember the RBR state.
	// tableMaps is indexed by tableID.
	tableMaps := make(map[uint64]*tableCacheEntry)

	// // A begin can be triggered either by a BEGIN query, or by a GTID_EVENT.
	sendMsg := func() {
		if autocommit {
			pbMsgCh <- &ParsedProto{
				ErrInfo: err,
				Events:  pbEvents,
				NextPos: pos,
				Reload:  reload,
			}
			// reset err and events
			err, pbEvents = nil, nil
			reload = false
		}
	}

	// Parse events.
	for {
		var ev mysql.BinlogEvent
		var ok bool

		select {
		case ev, ok = <-events:
			if !ok {
				// events channel has been closed, which means the connection died.
				// log.Infof("reached end of binlog event stream")
				err = ErrServerEOF
				sendMsg()
				break
			}
		case <-ctx.Done():
			err = ctx.Err()
			sendMsg()
			break
		}

		if ev == nil {
			continue
		}

		// Validate the buffer before reading fields from it.
		if !ev.IsValid() {
			err = fmt.Errorf("can't parse binlog event, invalid data: %#v", ev)
			sendMsg()
			continue
		}

		// We need to keep checking for FORMAT_DESCRIPTION_EVENT even after we've
		// seen one, because another one might come along (e.g. on log rotate due to
		// binlog settings change) that changes the format.
		if ev.IsFormatDescription() {
			format, err = ev.Format()
			if err != nil {
				err = fmt.Errorf("can't parse FORMAT_DESCRIPTION_EVENT: %v, event data: %#v", err, ev)
				sendMsg()
			}
			continue
		}

		// We can't parse anything until we get a FORMAT_DESCRIPTION_EVENT that
		// tells us the size of the event header.
		if format.IsZero() {
			// The only thing that should come before the FORMAT_DESCRIPTION_EVENT
			// is a fake ROTATE_EVENT, which the master sends to tell us the name
			// of the current log file.
			if !ev.IsRotate() {
				err = fmt.Errorf("got a real event before FORMAT_DESCRIPTION_EVENT: %#v", ev)
				sendMsg()
			}
			continue
		}
		// Strip the checksum, if any. We don't actually verify the checksum, so discard it.
		ev, _, err = ev.StripChecksum(format)
		if err != nil {
			err = fmt.Errorf("can't strip checksum from binlog event: %v, event data: %#v", err, ev)
			sendMsg()
			continue
		}
		switch {
		case ev.IsPseudo():
			gtid, _, err = ev.GTID(format)
			if err != nil {
				err = fmt.Errorf("can't get GTID from binlog event: %v, event data: %#v", err, ev)
				sendMsg()
				continue
			}
			oldpos := pos
			pos = mysql.AppendGTID(pos, gtid)
			// If the event is received outside of a transaction, it must
			// be sent. Otherwise, it will get lost and the targets will go out
			// of sync.
			if pos.Equal(oldpos) {
				err = fmt.Errorf("the same pos")
				sendMsg()
				continue
			}
		case ev.IsGTID(): // GTID_EVENT: update current GTID, maybe BEGIN.
			gtid, _, err = ev.GTID(format)
			if err != nil {
				err = fmt.Errorf("can't get GTID from binlog event: %v, event data: %#v", err, ev)
				sendMsg()
				continue
			}
			pos = mysql.AppendGTID(pos, gtid)
		case ev.IsXID(): // XID_EVENT (equivalent to COMMIT)
			autocommit = true
			sendMsg()
		case ev.IsIntVar(): // INTVAR_EVENT
			_, _, err = ev.IntVar(format)
			if err != nil {
				err = fmt.Errorf("can't parse INTVAR_EVENT: %v, event data: %#v", err, ev)
				sendMsg()
				continue
			}
			// todo deal with intvar type event
		case ev.IsRand(): // RAND_EVENT
			_, _, err = ev.Rand(format)
			if err != nil {
				err = fmt.Errorf("can't parse RAND_EVENT: %v, event data: %#v", err, ev)
				sendMsg()
				continue
			}
			// todo deal with random type event
		case ev.IsQuery(): // QUERY_EVENT
			// Extract the query string and group into transactions.
			q, err := ev.Query(format)
			if err != nil {
				err = fmt.Errorf("can't get query from binlog event: %v, event data: %#v", err, ev)
				sendMsg()
				continue
			}
			switch cat := getStatementCategory(q.SQL); cat {
			case binlogdatapb.BinlogTransaction_Statement_BL_BEGIN:
				autocommit = false
			case binlogdatapb.BinlogTransaction_Statement_BL_ROLLBACK:
				// Rollbacks are possible under some circumstances. Since the stream
				// client keeps track of its replication position by updating the set
				// of GTIDs it's seen, we must commit an empty transaction so the client
				// can update its position.
				pbEvents = nil
				autocommit = true
				fallthrough
			case binlogdatapb.BinlogTransaction_Statement_BL_COMMIT:
				autocommit = true
				sendMsg()
			case binlogdatapb.BinlogTransaction_Statement_BL_DDL:
				autocommit, reload = true, true
				sendMsg()
			default: // BL_DDL, BL_SET, BL_INSERT, BL_UPDATE, BL_DELETE, BL_UNRECOGNIZED
				// todo do actions here
			}
		case ev.IsPreviousGTIDs(): // PREVIOUS_GTIDS_EVENT
			// MySQL 5.6 only: The Binlogs contain an
			// event that gives us all the previously
			// applied commits. It is *not* an
			// authoritative value, unless we started from
			// the beginning of a binlog file.

			// set as start pos
			pos = startPos
		case ev.IsTableMap():
			// Save all tables, even not in the same DB.
			tableID := ev.TableID(format)
			var tm *mysql.TableMap
			tm, err = ev.TableMap(format)
			if err != nil {
				sendMsg()
				continue
			}
			// TODO(alainjobart) if table is already in map,
			// just use it.

			tce := &tableCacheEntry{
				tm:              tm,
				keyspaceIDIndex: -1,
			}
			tableMaps[tableID] = tce

			// Find and fill in the table schema.
			se, ok := seList[tm.Database]
			if !ok {
				err = fmt.Errorf("schema %v is not monitored", tm.Database)
				continue
			}
			tce.ti = se.GetTable(sqlparser.NewTableIdent(tm.Name))
			if tce.ti == nil {
				err = fmt.Errorf("unknown table %v in schema %v", tm.Name, tm.Database)
				continue
			}

			tce.pkNames = make([]*querypb.Field, len(tce.ti.PKColumns))
			tce.pkIndexes = make([]int, len(tce.ti.Columns))
			for i := range tce.pkIndexes {
				// Put -1 as default in here.
				tce.pkIndexes[i] = -1
			}
			for i, c := range tce.ti.PKColumns {
				// Patch in every PK column index.
				tce.pkIndexes[c] = i
				// Fill in pknames
				tce.pkNames[i] = &querypb.Field{
					Name: tce.ti.Columns[c].Name.String(),
					Type: tce.ti.Columns[c].Type,
				}
			}
		case ev.IsWriteRows():
			tableID := ev.TableID(format)
			tce, ok := tableMaps[tableID]
			if !ok {
				err = fmt.Errorf("unknown tableID %v in InsertRows event", tableID)
				continue
			}
			var rows mysql.Rows
			rows, err = ev.Rows(format, tce.tm)
			if err == nil {
				var insertEvents []*pbmysql.Event
				insertEvents, err = transToProto(tce, &rows, pbmysql.EventType_InsertEvent, onlyPrimaryID)
				if err == nil {
					pbEvents = append(pbEvents, insertEvents...)
				}
			}
			sendMsg()
		case ev.IsUpdateRows():
			tableID := ev.TableID(format)
			tce, ok := tableMaps[tableID]
			if !ok {
				err = fmt.Errorf("unknown tableID %v in UpdateRows event", tableID)
				sendMsg()
				continue
			}

			var rows mysql.Rows
			rows, err = ev.Rows(format, tce.tm)
			if err == nil {
				var updateEvent []*pbmysql.Event
				updateEvent, err = transToProto(tce, &rows, pbmysql.EventType_UpdateEvent, onlyPrimaryID)
				if err == nil {
					pbEvents = append(pbEvents, updateEvent...)
				}
			}
			// statements = appendUpdates(statements, tce, &rows)
			sendMsg()
		case ev.IsDeleteRows():
			tableID := ev.TableID(format)
			tce, ok := tableMaps[tableID]
			if !ok {
				err = fmt.Errorf("unknown tableID %v in DeleteRows event", tableID)
				sendMsg()
				continue
			}

			var rows mysql.Rows
			rows, err = ev.Rows(format, tce.tm)

			if err == nil {
				var deleteEvent []*pbmysql.Event
				deleteEvent, err = transToProto(tce, &rows, pbmysql.EventType_DeleteEvent, onlyPrimaryID)
				if err == nil {
					pbEvents = append(pbEvents, deleteEvent...)
				}
			}
			// statements = appendDeletes(statements, tce, &rows)
			sendMsg()
		}
	}
}

// transfer to proto type
func transToProto(tce *tableCacheEntry, rows *mysql.Rows, et pbmysql.EventType, onlyPrimaryID bool) (pbEvents []*pbmysql.Event, err error) {
	for i := range rows.Rows {
		// sql := sqlparser.NewTrackedBuffer(nil)
		e := &pbmysql.Event{
			Schema:  tce.tm.Database,
			Table:   tce.tm.Name,
			Columns: []*pbmysql.ColumnValue{},
			Et:      et,
		}
		e.Id, e.Columns, err = transToProtoColumn(tce, rows, i, et == pbmysql.EventType_DeleteEvent, onlyPrimaryID)
		if err != nil {
			return
		} else if onlyPrimaryID && e.Id == 0 {
			// id must be > 0`
			continue
		}
		pbEvents = append(pbEvents, e)
	}
	return
}

func transToProtoColumn(tce *tableCacheEntry, rs *mysql.Rows, rowIndex int, isDelete bool, onlyPrimaryID bool) (int64, []*pbmysql.ColumnValue, error) {
	var primaryID int64
	var columns []*pbmysql.ColumnValue

	// must be primary key with type
	if onlyPrimaryID {
		if len(tce.pkNames) != 1 {
			return 0, nil, nil
		}
	}

	valueIndex := 0
	pos := 0
	var data []byte
	var bitColumns mysql.Bitmap
	if isDelete {
		data = rs.Rows[rowIndex].Identify
		bitColumns = rs.IdentifyColumns
	} else {
		data = rs.Rows[rowIndex].Data
		bitColumns = rs.DataColumns
	}
	for c := 0; c < bitColumns.Count(); c++ {
		oneColumn := &pbmysql.ColumnValue{}
		if !bitColumns.Bit(c) {
			continue
		}
		// set column name
		oneColumn.Name = tce.ti.Columns[c].Name.String()

		if isDelete {
			if rs.Rows[rowIndex].NullIdentifyColumns.Bit(valueIndex) {
				// This column is represented, but its value is NULL.
				// oneColumn.Value = []byte("null")
				// columns = append(columns, oneColumn)
				valueIndex++
				continue
			}
		} else {
			if rs.Rows[rowIndex].NullColumns.Bit(valueIndex) {
				// This column is represented, but its value is NULL.
				// oneColumn.Value = []byte("null")
				// columns = append(columns, oneColumn)
				valueIndex++
				continue
			}
		}

		// We have real data.
		value, l, err := mysql.CellValue(data, pos, tce.tm.Types[c], tce.tm.Metadata[c], tce.ti.Columns[c].Type)
		if err != nil {
			return 0, nil, err
		}

		if oneColumn.Name == tce.pkNames[0].GetName() {
			primaryID, err = strconv.ParseInt(value.ToString(), 10, 64)
			if err != nil {
				return 0, nil, err
			}
		}
		// todo timestamp format will have problem
		// if value.Type() == querypb.Type_TIMESTAMP && !bytes.HasPrefix(value.ToBytes(), mysql.ZeroTimestamp) {
		// 	// Values in the binary log are UTC. Let's convert them
		// 	// to whatever timezone the connection is using,
		// 	// so MySQL properly converts them back to UTC.
		// 	sql.WriteString("convert_tz(")
		// 	value.EncodeSQL(sql)
		// 	sql.WriteString(", '+00:00', @@session.time_zone)")
		// } else {
		// value.EncodeSQL(sql)
		// }
		oneColumn.Value = value.ToBytes()

		pos += l
		valueIndex++

		columns = append(columns, oneColumn)
	}

	return primaryID, columns, nil
}
