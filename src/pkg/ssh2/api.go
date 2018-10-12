package ssh2

import (
	"time"
)

// SSHConfig for ssh to certain machine
type SSHConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Key      string `json:"key" yaml:"key"`
}

// Clienter : interface of ssh client
type Clienter interface {
	RunCmd(cmd string) (out string, err error)
	Monitor() (out string)
	GetStats() (stats *StatsInfo)
	ClearStats()
}

// NewClienter : implements ssh client to remote machines
func NewClienter(sshCfg *SSHConfig) (clnter Clienter, err error) {
	return sshCfg.newClient()
}

// FSInfo file system info
type FSInfo struct {
	MountPoint string
	Used       uint64
	Free       uint64
}

// NetIntfInfo net info
type NetIntfInfo struct {
	IPv4 string
	IPv6 string
	Rx   uint64
	Tx   uint64
}

// CPUInfo cpu info
type CPUInfo struct {
	User    float32
	Nice    float32
	System  float32
	Idle    float32
	Iowait  float32
	Irq     float32
	SoftIrq float32
	Steal   float32
	Guest   float32
}

// StatsInfo : information collected from machine
type StatsInfo struct {
	Uptime       time.Duration
	Hostname     string
	Load1        string
	Load5        string
	Load10       string
	RunningProcs string
	TotalProcs   string
	MemTotal     uint64
	MemFree      uint64
	MemBuffers   uint64
	MemCached    uint64
	SwapTotal    uint64
	SwapFree     uint64
	FSInfos      []FSInfo
	NetIntf      map[string]NetIntfInfo
	CPU          CPUInfo
	preCPU       cpuRaw
}
