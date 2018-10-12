package ssh2

import (
	"bytes"
	"encoding/json"
)

// format result like json string
func (clnt *sshClient) showStats() (output string) {
	clnt.getAllStats()
	b, _ := json.Marshal(clnt.Stats)
	var out bytes.Buffer
	json.Indent(&out, b, "", "    ")
	output = out.String()
	return
}
