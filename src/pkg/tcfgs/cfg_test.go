package tcfgs

import "testing"

func TestGetTestMysqlCfg(t *testing.T) {
	_, err := GetTestMysqlCfg()
	if err != nil {
		t.Errorf("read mysql cfg error %v", err)
	}
}

func TestGetTestSSHCfg(t *testing.T) {
	_, err := GetTestSSHCfg()
	if err != nil {
		t.Errorf("read ssh cfg error %v", err)
	}
}

func TestGetTestKafkaCfg(t *testing.T) {
	_, err := GetTestKafkaCfg()
	if err != nil {
		t.Errorf("read kafka cfg error %v", err)
	}
}
