package gotimer

// SimpleTimer implement a simple timer
type SimpleTimer interface {
	// Only send msg channel once
	SendMsgCh()
	CheckTimeout() (isTimeout bool)
}

// NewSimpleMsgCh : new a timple timer
// timeout is used as seconds
func NewSimpleMsgCh(timeout int) SimpleTimer {
	return newSimpleMsgCh(timeout)
}
