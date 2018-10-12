package tcfgs

import (
	"pkg/cfgtool"
)

const cfgPath = "../../configs/testing/test.yaml"

var cfgsAll = &baseCfg{
	Mysql:   &MysqlCfg{},
	SSHHost: &SSHCfg{},
}

// GetTestMysqlCfg get test mysql config from yaml
func GetTestMysqlCfg() (*MysqlCfg, error) {
	err := cfgtool.LoadCfgFromYamlFile(cfgPath, cfgsAll)
	return cfgsAll.Mysql, err
}

// GetTestSSHCfg get test ssh config from yaml
func GetTestSSHCfg() (*SSHCfg, error) {
	err := cfgtool.LoadCfgFromYamlFile(cfgPath, cfgsAll)
	return cfgsAll.SSHHost, err
}

// MysqlCfg config for a mysql instance
type MysqlCfg struct {
	IP       string
	Port     int
	User     string
	Password string
}

// SSHCfg config for a mysql instance
type SSHCfg struct {
	Host     string
	Port     int
	Username string
	Password string
}

type baseCfg struct {
	Mysql   *MysqlCfg `yaml:"mysql"`
	SSHHost *SSHCfg   `yaml:"sshHost"`
}
