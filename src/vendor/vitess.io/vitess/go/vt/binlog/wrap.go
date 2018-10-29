package binlog

import (
	"fmt"
	"time"

	"vitess.io/vitess/go/vt/vttablet/tabletserver/schema"

	"golang.org/x/net/context"
	"vitess.io/vitess/go/mysql"
	binlogdatapb "vitess.io/vitess/go/vt/proto/binlogdata"
	querypb "vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/sqlparser"
)

// wrapped from binlog_streamer.go

// ParsedStat used to parse event to statmentss
type ParsedStat struct {
	NextPos mysql.Position
	ErrInfo error
	Stats   []FullBinlogStatement
}

// ParseStatsEvents used to parse event detail from binlog
// onlyPrimaryID if is true, filter records without the primary key named "id"
func ParseStatsEvents(ctx context.Context, events <-chan mysql.BinlogEvent, seList map[string]*schema.Engine, pbMsgCh chan *ParsedStat) {
	// statements is for sql statement
	var statements []FullBinlogStatement
	var format mysql.BinlogFormat
	var gtid mysql.GTID
	var pos mysql.Position
	var autocommit = true
	var err error

	// Remember the RBR state.
	// tableMaps is indexed by tableID.
	tableMaps := make(map[uint64]*tableCacheEntry)

	// // A begin can be triggered either by a BEGIN query, or by a GTID_EVENT.
	sendMsg := func() {
		if autocommit {
			pbMsgCh <- &ParsedStat{
				ErrInfo: err,
				Stats:   statements,
				NextPos: pos,
			}
			// reset err and statements
			err, statements = nil, nil
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
			typ, value, err := ev.IntVar(format)
			if err != nil {
				err = fmt.Errorf("can't parse INTVAR_EVENT: %v, event data: %#v", err, ev)
				sendMsg()
				continue
			}
			statements = append(statements, FullBinlogStatement{
				Statement: &binlogdatapb.BinlogTransaction_Statement{
					Category: binlogdatapb.BinlogTransaction_Statement_BL_SET,
					Sql:      []byte(fmt.Sprintf("SET %s=%d", mysql.IntVarNames[typ], value)),
				},
			})
			sendMsg()
		case ev.IsRand(): // RAND_EVENT
			fmt.Println("Is rand")
			seed1, seed2, err := ev.Rand(format)
			if err != nil {
				err = fmt.Errorf("can't parse RAND_EVENT: %v, event data: %#v", err, ev)
				sendMsg()
				continue
			}
			statements = append(statements, FullBinlogStatement{
				Statement: &binlogdatapb.BinlogTransaction_Statement{
					Category: binlogdatapb.BinlogTransaction_Statement_BL_SET,
					Sql:      []byte(fmt.Sprintf("SET @@RAND_SEED1=%d, @@RAND_SEED2=%d", seed1, seed2)),
				},
			})
			sendMsg()
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
				statements = nil
				autocommit = true
				fallthrough
			case binlogdatapb.BinlogTransaction_Statement_BL_COMMIT:
				fmt.Println("Is commit")
				autocommit = true
				sendMsg()
			default: // BL_DDL, BL_SET, BL_INSERT, BL_UPDATE, BL_DELETE, BL_UNRECOGNIZED
				setTimestamp := &binlogdatapb.BinlogTransaction_Statement{
					Category: binlogdatapb.BinlogTransaction_Statement_BL_SET,
					Sql:      []byte(fmt.Sprintf("SET TIMESTAMP=%d", ev.Timestamp())),
				}
				statement := &binlogdatapb.BinlogTransaction_Statement{
					Category: cat,
					Sql:      []byte(q.SQL),
				}
				statements = append(statements, FullBinlogStatement{
					Statement: setTimestamp,
				}, FullBinlogStatement{
					Statement: statement,
				})
				sendMsg()
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
			fmt.Println("Is write ,table id is", tableID)
			tce, ok := tableMaps[tableID]
			if !ok {
				err = fmt.Errorf("unknown tableID %v in UpdateRows event", tableID)
				continue
			}
			var rows mysql.Rows
			rows, err = ev.Rows(format, tce.tm)
			if err == nil {
				statements = appendInserts(statements, tce, &rows)
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
				statements = appendUpdates(statements, tce, &rows)
			}
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
				statements = appendDeletes(statements, tce, &rows)
			}
			sendMsg()
		}
	}
}

func appendInserts(statements []FullBinlogStatement, tce *tableCacheEntry, rows *mysql.Rows) []FullBinlogStatement {
	bls := NewStreamer(nil, nil, nil, mysql.Position{}, time.Now().Unix(), nil)
	return bls.appendInserts(statements, tce, rows)
}

func appendUpdates(statements []FullBinlogStatement, tce *tableCacheEntry, rows *mysql.Rows) []FullBinlogStatement {
	bls := NewStreamer(nil, nil, nil, mysql.Position{}, time.Now().Unix(), nil)
	return bls.appendUpdates(statements, tce, rows)
}

func appendDeletes(statements []FullBinlogStatement, tce *tableCacheEntry, rows *mysql.Rows) []FullBinlogStatement {
	bls := NewStreamer(nil, nil, nil, mysql.Position{}, time.Now().Unix(), nil)
	return bls.appendDeletes(statements, tce, rows)
}
