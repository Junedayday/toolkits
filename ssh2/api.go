package ssh2

// ConnCfger : implement the configuration of a ssh connection
type ConnCfger interface {
	NewCmdClienter() (clnter CmdClienter, err error)
	NewMachineStatsClienter() (clnter MachineStatsClienter, err error)
}

// NewConnCfger : user,password,ip,port are essential info
func NewConnCfger(user, password, ip string, port int) ConnCfger {
	return newConnCfg(user, password, ip, port)
}

// CmdClienter : implemnent for run cmd by ssh
type CmdClienter interface {
	RunCmd(cmd string) (output string, err error)
}

// NewCmdClienter : create a ssh client to run cmds
func (cfg *connCfg) NewCmdClienter() (clnter CmdClienter, err error) {
	return cfg.newClient()
}

// MachineStatsClienter : implemnent for get machine stats by ssh
type MachineStatsClienter interface {
	GetMachineStats() (stats string)
}

// NewMachineStatsClienter : create a ssh client to collecting machine stats
func (cfg *connCfg) NewMachineStatsClienter() (clnter MachineStatsClienter, err error) {
	return cfg.newClient()
}
