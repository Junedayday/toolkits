package ssh2

import (
	"bytes"
	"pkg/gotimer"

	"golang.org/x/crypto/ssh"
)

const (
	defaultSessionTimeout = 3
	defaultCmdTimeout     = 3
)

type sshClient struct {
	*ssh.Client
	Stats *StatsInfo
}

func (clnt *sshClient) RunCmd(cmd string) (out string, err error) {
	var session *ssh.Session
	// set timeout for session
	session, err = clnt.newSessionWithTimeout(defaultSessionTimeout)
	if err != nil {
		return
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	session.Stderr = &buf
	// set timeout for run command here
	err = runCmdWithTimeout(session, cmd, defaultCmdTimeout)
	newByte := getValidByte([]byte(buf.String()))
	out = string(newByte)
	session.Close()
	return
}

func (clnt *sshClient) Monitor() (out string) {
	out = clnt.showStats()
	return
}

func (clnt *sshClient) GetStats() (stats *StatsInfo) {
	stats = clnt.Stats
	return
}

func (clnt *sshClient) ClearStats() {
	clnt.clearStat()
	return
}

func (clnt *sshClient) newSessionWithTimeout(interval int) (session *ssh.Session, err error) {
	st := gotimer.NewSimpleMsgCh(interval)
	go func() {
		session, err = clnt.NewSession()
		st.SendMsgCh()
	}()

	if st.CheckTimeout() {
		return nil, errSessionTimeout
	}
	return
}

func runCmdWithTimeout(session *ssh.Session, cmd string, interval int) (err error) {
	st := gotimer.NewSimpleMsgCh(interval)
	go func() {
		err = session.Run(cmd)
		st.SendMsgCh()
	}()

	if st.CheckTimeout() {
		return errCmdTimeout
	}
	return
}

// remove invalid byte of 0
func getValidByte(src []byte) []byte {
	var strBuf []byte
	for _, v := range src {
		if v != 0 {
			strBuf = append(strBuf, v)
		}
	}
	return strBuf
}
