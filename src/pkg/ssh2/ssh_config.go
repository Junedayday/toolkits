package ssh2

import (
	"fmt"
	"net"
	"pkg/gotimer"

	"golang.org/x/crypto/ssh"
)

const (
	clntDefaultTimeout = 3
)

// parse public key from private key
func (sshCfg *SSHConfig) newClient() (clnt *sshClient, err error) {
	config := &ssh.ClientConfig{
		User: sshCfg.Username,
	}
	if len(sshCfg.Password) == 0 {
		var key ssh.Signer
		key, err = ssh.ParsePrivateKey([]byte(sshCfg.Key))
		if err != nil {
			return
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(key)}
	} else {
		config.Auth = []ssh.AuthMethod{ssh.Password(sshCfg.Password)}
		config.HostKeyCallback = func(string, net.Addr, ssh.PublicKey) error { return nil }
	}
	var client *ssh.Client

	st := gotimer.NewSimpleMsgCh(clntDefaultTimeout)
	go func() {
		client, err = ssh.Dial("tcp", fmt.Sprintf("%v:%v", sshCfg.Host, sshCfg.Port), config)
		st.SendMsgCh()
	}()

	if st.CheckTimeout() {
		return nil, errClientTimeout
	} else if err == nil {
		clnt = &sshClient{
			Client: client,
			Stats: &StatsInfo{
				NetIntf: make(map[string]NetIntfInfo),
			},
		}
	}
	return
}
