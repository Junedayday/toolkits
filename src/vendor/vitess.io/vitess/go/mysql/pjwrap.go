package mysql

import (
	"encoding/binary"
	"fmt"
)

func (ev binlogEvent) NextPosition() int64 {
	return int64(binary.LittleEndian.Uint32(ev.Bytes()[13 : 13+4]))
}

// Format implements BinlogEvent.Format().
//
// Expected rotate (L = total length of event data):
//   # bytes    field
//     8         position
//    L-8        filename
func (ev binlogEvent) Rotate(f BinlogFormat) (string, int64, error) {
	data := ev.Bytes()[f.HeaderLength:]
	if len(data) < 8 {
		return "", 0, fmt.Errorf("Rotate position overflows buffer (8 > %v)", len(data))
	}
	offset := int64(binary.LittleEndian.Uint64(data[0:8]))
	filename := string(data[8:])
	return filename, offset, nil
}
