package cfgtool

import (
	"testing"
)

func TestLoadYamlFile(t *testing.T) {
	type MysqlCfg struct {
		IP       string
		Port     int
		User     string
		Password string
	}
	type yamlCfg struct {
		Mysql MysqlCfg
	}
	testCfg := &yamlCfg{}
	err := LoadCfgFromYamlFile("../../configs/testing/test.yaml", testCfg)
	if err != nil {
		t.Error("Load yaml file failed!")
		return
	}
	if testCfg.Mysql.IP == "" || testCfg.Mysql.Port == 0 || testCfg.Mysql.User == "" || testCfg.Mysql.Password == "" {
		t.Error("read yaml file wrong!")
		return
	}

	// defined a wrong path
	err = LoadCfgFromYamlFile("../../configs/testing/testing1.yaml", testCfg)
	if err == nil {
		t.Error("file not exist")
	}
}
