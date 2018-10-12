package binlog

import (
	"fmt"
	"proto/mysql/pbmysql"
	"strconv"

	"vitess.io/vitess/go/sqltypes"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/schema"

	"golang.org/x/net/context"
	"vitess.io/vitess/go/mysql"
	binlogdatapb "vitess.io/vitess/go/vt/proto/binlogdata"
	"vitess.io/vitess/go/vt/sqlparser"
)

// wrapped from binlog_streamer.go

// ParsedProto used to parse event to proto
type ParsedProto struct {
	NextPos mysql.Position
	ErrInfo error
	Events  []*pbmysql.Event
}

// ParseProtoEvents used to parse event detail from binlog
func ParseProtoEvents(ctx context.Context, events <-chan mysql.BinlogEvent, seList map[string]*schema.Engine, pbMsgCh chan *ParsedProto) {
	var format mysql.BinlogFormat
	var gtid mysql.GTID
	var pos mysql.Position
	var autocommit = true
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
			}
			// reset err and events
			err, pbEvents = nil, nil
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
			fmt.Println("Is xid")
			autocommit = true
			sendMsg()
		case ev.IsIntVar(): // INTVAR_EVENT
			fmt.Println("Is int")
			_, _, err = ev.IntVar(format)
			if err != nil {
				err = fmt.Errorf("can't parse INTVAR_EVENT: %v, event data: %#v", err, ev)
				sendMsg()
				continue
			}
			// todo deal with intvar type event
		case ev.IsRand(): // RAND_EVENT
			fmt.Println("Is rand")
			_, _, err = ev.Rand(format)
			if err != nil {
				err = fmt.Errorf("can't parse RAND_EVENT: %v, event data: %#v", err, ev)
				sendMsg()
				continue
			}
			// todo deal with random type event
		case ev.IsQuery(): // QUERY_EVENT
			fmt.Println("Is query")
			// Extract the query string and group into transactions.
			q, err := ev.Query(format)
			if err != nil {
				err = fmt.Errorf("can't get query from binlog event: %v, event data: %#v", err, ev)
				sendMsg()
				continue
			}
			switch cat := getStatementCategory(q.SQL); cat {
			case binlogdatapb.BinlogTransaction_Statement_BL_BEGIN:
				fmt.Println("Is begin")
				autocommit = false
			case binlogdatapb.BinlogTransaction_Statement_BL_ROLLBACK:
				// Rollbacks are possible under some circumstances. Since the stream
				// client keeps track of its replication position by updating the set
				// of GTIDs it's seen, we must commit an empty transaction so the client
				// can update its position.
				fmt.Println("Is roll back")
				pbEvents = nil
				autocommit = true
				fallthrough
			case binlogdatapb.BinlogTransaction_Statement_BL_COMMIT:
				fmt.Println("Is commit")
				autocommit = true
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
			var newPos mysql.Position
			newPos, err = ev.PreviousGTIDs(format)
			if err != nil {
				sendMsg()
				continue
			}
			pos = newPos
		case ev.IsTableMap():
			// Save all tables, even not in the same DB.
			tableID := ev.TableID(format)
			fmt.Println("Is table map ,table id is", tableID)
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
		case ev.IsWriteRows():
			tableID := ev.TableID(format)
			fmt.Println("Is write ,table id is", tableID)
			tce, ok := tableMaps[tableID]
			if !ok {
				err = fmt.Errorf("unknown tableID %v in UpdateRows event", tableID)
				continue
			}
			var rows mysql.Rows
			rows, err = ev.Rows(format, tce.tm)
			if err == nil {
				var insertEvents []*pbmysql.Event
				insertEvents, err = transToProto(tce, &rows, pbmysql.EventType_InsertEvent)
				if err == nil {
					pbEvents = append(pbEvents, insertEvents...)
				}
			}
			sendMsg()
		case ev.IsUpdateRows():
			tableID := ev.TableID(format)
			fmt.Println("Is update,table id is", tableID)
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
				updateEvent, err = transToProto(tce, &rows, pbmysql.EventType_UpdateEvent)
				if err == nil {
					pbEvents = append(pbEvents, updateEvent...)
				}
			}
			// statements = appendUpdates(statements, tce, &rows)

			sendMsg()
		case ev.IsDeleteRows():
			tableID := ev.TableID(format)
			fmt.Println("Is delete,table id is", tableID)
			tce, ok := tableMaps[tableID]
			if !ok {
				err = fmt.Errorf("unknown tableID %v in UpdateRows event", tableID)
				sendMsg()
				continue
			}

			var rows mysql.Rows
			rows, err = ev.Rows(format, tce.tm)
			if err == nil {
				var deleteEvent []*pbmysql.Event
				deleteEvent, err = transToProto(tce, &rows, pbmysql.EventType_DeleteEvent)
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
func transToProto(tce *tableCacheEntry, rows *mysql.Rows, et pbmysql.EventType) (pbEvents []*pbmysql.Event, err error) {
	for i := range rows.Rows {
		// sql := sqlparser.NewTrackedBuffer(nil)
		e := &pbmysql.Event{
			Schema:  tce.tm.Database,
			Table:   tce.tm.Name,
			Columns: []*pbmysql.ColumnValue{},
			Et:      et,
		}
		e.Id, e.Columns, err = transToProtoColumn(tce, rows, i, tce.pkNames != nil)
		if err != nil {
			return
		}
		pbEvents = append(pbEvents, e)
	}
	return
}

func transToProtoColumn(tce *tableCacheEntry, rs *mysql.Rows, rowIndex int, getPK bool) (int64, []*pbmysql.ColumnValue, error) {
	var primaryID int64
	var columns []*pbmysql.ColumnValue

	valueIndex := 0
	data := rs.Rows[rowIndex].Data
	pos := 0
	//var keyspaceIDCell sqltypes.Value
	var pkValues []sqltypes.Value
	if getPK {
		pkValues = make([]sqltypes.Value, len(tce.pkNames))
	}
	for c := 0; c < rs.DataColumns.Count(); c++ {
		oneColumn := &pbmysql.ColumnValue{}
		if !rs.DataColumns.Bit(c) {
			continue
		}
		// set column name
		oneColumn.ColumnName = tce.ti.Columns[c].Name.String()

		if rs.Rows[rowIndex].NullColumns.Bit(valueIndex) {
			// This column is represented, but its value is NULL.
			oneColumn.Value = []byte("null")
			valueIndex++
			continue
		}

		// We have real data.
		value, l, err := mysql.CellValue(data, pos, tce.tm.Types[c], tce.tm.Metadata[c], tce.ti.Columns[c].Type)
		if err != nil {
			return 0, nil, err
		}

		if oneColumn.ColumnName == "id" {
			primaryID, err = strconv.ParseInt(value.ToString(), 10, 64)
			if err != nil {
				return 0, nil, err
			}
			pos += l
			valueIndex++
			continue
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
		// if c == tce.keyspaceIDIndex {
		// 	keyspaceIDCell = value
		// }
		if getPK {
			if tce.pkIndexes[c] != -1 {
				pkValues[tce.pkIndexes[c]] = value
			}
		}
		pos += l
		valueIndex++

		columns = append(columns, oneColumn)
	}

	return primaryID, columns, nil
}
