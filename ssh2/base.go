package ssh2

import (
	"bytes"
	"encoding/json"

	"golang.org/x/crypto/ssh"
)

type client struct {
	*ssh.Client
	stats *statsInfo
}

func (clnt *client) RunCmd(cmd string) (output string, err error) {
	var session *ssh.Session
	if session, err = clnt.NewSession(); err != nil {
		return
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	session.Stderr = &buf
	err = session.Run(cmd)
	output = string(getValidByte([]byte(buf.String())))
	return
}

func (clnt *client) GetMachineStats() (stats string) {
	stats = clnt.getStatsByString()
	return
}

func (clnt *client) getStatsByString() (stats string) {
	clnt.getAllStats()
	b, _ := json.Marshal(clnt.stats)
	var buf bytes.Buffer
	json.Indent(&buf, b, "", "    ")
	stats = buf.String()
	return
}

// remove invalid byte of 0
func getValidByte(src []byte) (dest []byte) {
	for _, v := range src {
		if v != 0 {
			dest = append(dest, v)
		}
	}
	return
}
