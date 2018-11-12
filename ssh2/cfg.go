package ssh2

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

const (
	errSSHNoPassword = "SSH only support password mode now"
)

type connCfg struct {
	User     string
	Password string
	IP       string
	Port     int
}

func newConnCfg(user, password, ip string, port int) (cfg *connCfg) {
	cfg = &connCfg{
		IP:       ip,
		Port:     port,
		User:     user,
		Password: password,
	}
	return
}

func (cfg *connCfg) newClient() (clnt *client, err error) {
	config := &ssh.ClientConfig{User: cfg.User}
	if len(cfg.Password) == 0 {
		err = fmt.Errorf(errSSHNoPassword)
		return
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(cfg.Password)}
	var sshClnt *ssh.Client
	sshClnt, err = ssh.Dial("tcp", fmt.Sprintf("%v:%v", cfg.IP, cfg.Port), config)
	clnt = &client{
		Client: sshClnt,
		stats: &statsInfo{
			NetInfos: make(map[string]netInfo),
		},
	}
	return
}
