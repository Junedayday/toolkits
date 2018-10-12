package gotimer

import (
	"testing"
	"time"
)

func TestCheckTimeout(t *testing.T) {
	st1 := NewSimpleMsgCh(2)
	go func() {
		time.Sleep(10 * time.Millisecond)
		st1.SendMsgCh()
	}()
	if st1.CheckTimeout() {
		t.Errorf("error time out")
	}

	st2 := NewSimpleMsgCh(1)
	go func() {
		time.Sleep(3 * time.Second)
		st2.SendMsgCh()
	}()
	if !st2.CheckTimeout() {
		t.Errorf("supposed to be time out")
	}
}
