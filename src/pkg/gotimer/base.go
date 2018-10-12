package gotimer

import (
	"time"
)

type simpleTimer struct {
	msgCh chan struct{}
	tm    int
}

func newSimpleMsgCh(timeout int) *simpleTimer {
	return &simpleTimer{
		msgCh: make(chan struct{}, 1),
		tm:    timeout,
	}
}

func (st *simpleTimer) SendMsgCh() {
	st.msgCh <- struct{}{}
}

func (st *simpleTimer) CheckTimeout() (isTimeout bool) {
	tm := time.NewTimer(time.Duration(st.tm) * time.Second)
	for {
		select {
		case <-st.msgCh:
			return
		case <-tm.C:
			isTimeout = true
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
