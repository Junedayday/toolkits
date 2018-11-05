package ssh2

import (
	"bufio"
	"strconv"
	"strings"
	"time"
)

type fsInfo struct {
	MountPoint string
	Used       uint64
	Free       uint64
}

type netInfo struct {
	IPv4 string
	IPv6 string
	Rx   uint64
	Tx   uint64
}

type cpuInfo struct {
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

type statsInfo struct {
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
	FSInfos      []fsInfo
	NetInfos     map[string]netInfo
	CPU          cpuInfo
	preCPU       cpuRaw
}

// TODO: this part is copied from web
// may not be so accurate for different systems
func (clnt *client) getAllStats() {
	clnt.getUptime()
	clnt.getHostname()
	clnt.getLoad()
	clnt.getMemInfo()
	clnt.getFSInfo()
	clnt.getInterfaces()
	clnt.getInterfaceInfo()
	clnt.getCPU()
}

func (clnt *client) getUptime() (err error) {
	uptime, err := clnt.RunCmd("/bin/cat /proc/uptime")
	if err != nil {
		return
	}

	parts := strings.Fields(uptime)
	if len(parts) == 2 {
		var upsecs float64
		upsecs, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return
		}
		clnt.stats.Uptime = time.Duration(upsecs * 1e9)
	}

	return
}

func (clnt *client) getHostname() (err error) {
	hostname, err := clnt.RunCmd("/bin/hostname -f")
	if err != nil {
		return
	}

	clnt.stats.Hostname = strings.TrimSpace(hostname)
	return
}

func (clnt *client) getLoad() (err error) {
	line, err := clnt.RunCmd("/bin/cat /proc/loadavg")
	if err != nil {
		return
	}

	parts := strings.Fields(line)
	if len(parts) == 5 {
		clnt.stats.Load1 = parts[0]
		clnt.stats.Load5 = parts[1]
		clnt.stats.Load10 = parts[2]
		if i := strings.Index(parts[3], "/"); i != -1 {
			clnt.stats.RunningProcs = parts[3][0:i]
			if i+1 < len(parts[3]) {
				clnt.stats.TotalProcs = parts[3][i+1:]
			}
		}
	}

	return
}

func (clnt *client) getMemInfo() (err error) {
	lines, err := clnt.RunCmd("/bin/cat /proc/meminfo")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 3 {
			val, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				continue
			}
			val *= 1024
			switch parts[0] {
			case "MemTotal:":
				clnt.stats.MemTotal = val
			case "MemFree:":
				clnt.stats.MemFree = val
			case "Buffers:":
				clnt.stats.MemBuffers = val
			case "Cached:":
				clnt.stats.MemCached = val
			case "SwapTotal:":
				clnt.stats.SwapTotal = val
			case "SwapFree:":
				clnt.stats.SwapFree = val
			}
		}
	}

	return
}

func (clnt *client) getFSInfo() (err error) {
	lines, err := clnt.RunCmd("/bin/df -B1")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	flag := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		n := len(parts)
		dev := n > 0 && strings.Index(parts[0], "/dev/") == 0
		if n == 1 && dev {
			flag = 1
		} else if (n == 5 && flag == 1) || (n == 6 && dev) {
			i := flag
			flag = 0
			used, err := strconv.ParseUint(parts[2-i], 10, 64)
			if err != nil {
				continue
			}
			free, err := strconv.ParseUint(parts[3-i], 10, 64)
			if err != nil {
				continue
			}
			clnt.stats.FSInfos = append(clnt.stats.FSInfos, fsInfo{
				parts[5-i], used, free,
			})
		}
	}

	return
}

func (clnt *client) getInterfaces() (err error) {
	var lines string
	lines, err = clnt.RunCmd("/bin/ip -o addr")
	if err != nil {
		// try /sbin/ip
		lines, err = clnt.RunCmd("/sbin/ip -o addr")
		if err != nil {
			return
		}
	}

	if clnt.stats.NetInfos == nil {
		clnt.stats.NetInfos = make(map[string]netInfo)
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 4 && (parts[2] == "inet" || parts[2] == "inet6") {
			ipv4 := parts[2] == "inet"
			intfname := parts[1]
			if info, ok := clnt.stats.NetInfos[intfname]; ok {
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				clnt.stats.NetInfos[intfname] = info
			} else {
				info := netInfo{}
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				clnt.stats.NetInfos[intfname] = info
			}
		}
	}

	return
}

func (clnt *client) getInterfaceInfo() (err error) {
	lines, err := clnt.RunCmd("/bin/cat /proc/net/dev")
	if err != nil {
		return
	}

	if clnt.stats.NetInfos == nil {
		return
	} // should have been here already

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 17 {
			intf := strings.TrimSpace(parts[0])
			intf = strings.TrimSuffix(intf, ":")
			if info, ok := clnt.stats.NetInfos[intf]; ok {
				rx, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					continue
				}
				tx, err := strconv.ParseUint(parts[9], 10, 64)
				if err != nil {
					continue
				}
				info.Rx = rx
				info.Tx = tx
				clnt.stats.NetInfos[intf] = info
			}
		}
	}

	return
}

type cpuRaw struct {
	User    uint64 // time spent in user mode
	Nice    uint64 // time spent in user mode with low priority (nice)
	System  uint64 // time spent in system mode
	Idle    uint64 // time spent in the idle task
	Iowait  uint64 // time spent waiting for I/O to complete (since Linux 2.5.41)
	Irq     uint64 // time spent servicing  interrupts  (since  2.6.0-test4)
	SoftIrq uint64 // time spent servicing softirqs (since 2.6.0-test4)
	Steal   uint64 // time spent in other OSes when running in a virtualized environment
	Guest   uint64 // time spent running a virtual CPU for guest operating systems under the control of the Linux kernel.
	Total   uint64 // total of all time fields
}

func (clnt *client) getCPU() (err error) {
	lines, err := clnt.RunCmd("/bin/cat /proc/stat")
	if err != nil {
		return
	}

	var (
		nowCPU cpuRaw
		total  float32
	)

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "cpu" { // changing here if want to get every cpu-core's stats
			parseCPUFields(fields, &nowCPU)
			break
		}
	}
	if clnt.stats.preCPU.Total == 0 { // having no pre raw cpu data
		goto END
	}

	total = float32(nowCPU.Total - clnt.stats.preCPU.Total)
	clnt.stats.CPU.User = float32(nowCPU.User-clnt.stats.preCPU.User) / total * 100
	clnt.stats.CPU.Nice = float32(nowCPU.Nice-clnt.stats.preCPU.Nice) / total * 100
	clnt.stats.CPU.System = float32(nowCPU.System-clnt.stats.preCPU.System) / total * 100
	clnt.stats.CPU.Idle = float32(nowCPU.Idle-clnt.stats.preCPU.Idle) / total * 100
	clnt.stats.CPU.Iowait = float32(nowCPU.Iowait-clnt.stats.preCPU.Iowait) / total * 100
	clnt.stats.CPU.Irq = float32(nowCPU.Irq-clnt.stats.preCPU.Irq) / total * 100
	clnt.stats.CPU.SoftIrq = float32(nowCPU.SoftIrq-clnt.stats.preCPU.SoftIrq) / total * 100
	clnt.stats.CPU.Guest = float32(nowCPU.Guest-clnt.stats.preCPU.Guest) / total * 100
END:
	clnt.stats.preCPU = nowCPU
	return
}

func parseCPUFields(fields []string, stat *cpuRaw) {
	numFields := len(fields)
	for i := 1; i < numFields; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}

		stat.Total += val
		switch i {
		case 1:
			stat.User = val
		case 2:
			stat.Nice = val
		case 3:
			stat.System = val
		case 4:
			stat.Idle = val
		case 5:
			stat.Iowait = val
		case 6:
			stat.Irq = val
		case 7:
			stat.SoftIrq = val
		case 8:
			stat.Steal = val
		case 9:
			stat.Guest = val
		}
	}
}
