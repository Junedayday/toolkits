package ssh2

import (
	"pkg/tcfgs"
	"testing"
)

// defined a right configuration
func testLoadSSHCfg() (*SSHConfig, error) {
	sshCfg, err := tcfgs.GetTestSSHCfg()
	if err != nil {
		return nil, err
	}
	return &SSHConfig{
		Host:     sshCfg.Host,
		Port:     sshCfg.Port,
		Username: sshCfg.Username,
		Password: sshCfg.Password,
	}, nil
}

func TestSSHClient(t *testing.T) {
	sshCfg1 := new(SSHConfig)
	_, err := NewClienter(sshCfg1)
	if err == nil {
		t.Errorf("must have password")
		return
	}

	sshCfg1.Password = "123456"
	_, err = NewClienter(sshCfg1)
	if err == nil {
		t.Errorf("connect to unkown server error %v", err)
		return
	}

	sshCfg, err := testLoadSSHCfg()
	if err != nil {
		t.Errorf("get ssh config error %v", err)
		return
	}
	_, err = NewClienter(sshCfg)
	if err != nil {
		t.Errorf("connect to server error : %v", err)
		return
	}
}

func TestSSHRunCmd(t *testing.T) {
	sshCfg, err := testLoadSSHCfg()
	if err != nil {
		t.Errorf("get ssh config error %v", err)
		return
	}
	clnt, err := NewClienter(sshCfg)
	if err != nil {
		t.Errorf("new ssh client error %v", err)
		return
	}
	out, err := clnt.RunCmd("ps")
	if err != nil {
		t.Errorf("run cmd error %v", err)
	} else if len(out) == 0 {
		t.Errorf("command has no return")
	}
}

func TestSSHStat(t *testing.T) {
	sshCfg, err := testLoadSSHCfg()
	if err != nil {
		t.Errorf("get ssh config error %v", err)
		return
	}
	clnt, err := NewClienter(sshCfg)
	if err != nil {
		t.Errorf("new ssh client error %v", err)
		return
	}

	out := clnt.Monitor()
	if len(out) == 0 {
		t.Errorf("no output")
		return
	}
}

func TestSSHGetStat(t *testing.T) {
	sshCfg, err := testLoadSSHCfg()
	if err != nil {
		t.Errorf("get ssh config error %v", err)
		return
	}
	clnt, err := NewClienter(sshCfg)
	if err != nil {
		t.Errorf("new ssh client error %v", err)
		return
	}
	out := clnt.Monitor()
	if len(out) == 0 {
		t.Errorf("no output")
		return
	}

	out2 := clnt.Monitor()
	if len(out2) == 0 {
		t.Errorf("no output2")
	}
}
