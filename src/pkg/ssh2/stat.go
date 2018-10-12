package ssh2

import (
	"bufio"
	"strconv"
	"strings"
	"time"
)

// this part is copied from web
// may not be so accurate for different systems

func (clnt *sshClient) getAllStats() {
	clnt.getUptime()
	clnt.getHostname()
	clnt.getLoad()
	clnt.getMemInfo()
	clnt.getFSInfo()
	clnt.getInterfaces()
	clnt.getInterfaceInfo()
	clnt.getCPU()
}

func (clnt *sshClient) clearStat() {
	clnt.Stats = &StatsInfo{
		NetIntf: make(map[string]NetIntfInfo),
	}
}

func (clnt *sshClient) getUptime() (err error) {
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
		clnt.Stats.Uptime = time.Duration(upsecs * 1e9)
	}

	return
}

func (clnt *sshClient) getHostname() (err error) {
	hostname, err := clnt.RunCmd("/bin/hostname -f")
	if err != nil {
		return
	}

	clnt.Stats.Hostname = strings.TrimSpace(hostname)
	return
}

func (clnt *sshClient) getLoad() (err error) {
	line, err := clnt.RunCmd("/bin/cat /proc/loadavg")
	if err != nil {
		return
	}

	parts := strings.Fields(line)
	if len(parts) == 5 {
		clnt.Stats.Load1 = parts[0]
		clnt.Stats.Load5 = parts[1]
		clnt.Stats.Load10 = parts[2]
		if i := strings.Index(parts[3], "/"); i != -1 {
			clnt.Stats.RunningProcs = parts[3][0:i]
			if i+1 < len(parts[3]) {
				clnt.Stats.TotalProcs = parts[3][i+1:]
			}
		}
	}

	return
}

func (clnt *sshClient) getMemInfo() (err error) {
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
				clnt.Stats.MemTotal = val
			case "MemFree:":
				clnt.Stats.MemFree = val
			case "Buffers:":
				clnt.Stats.MemBuffers = val
			case "Cached:":
				clnt.Stats.MemCached = val
			case "SwapTotal:":
				clnt.Stats.SwapTotal = val
			case "SwapFree:":
				clnt.Stats.SwapFree = val
			}
		}
	}

	return
}

func (clnt *sshClient) getFSInfo() (err error) {
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
			clnt.Stats.FSInfos = append(clnt.Stats.FSInfos, FSInfo{
				parts[5-i], used, free,
			})
		}
	}

	return
}

func (clnt *sshClient) getInterfaces() (err error) {
	var lines string
	lines, err = clnt.RunCmd("/bin/ip -o addr")
	if err != nil {
		// try /sbin/ip
		lines, err = clnt.RunCmd("/sbin/ip -o addr")
		if err != nil {
			return
		}
	}

	if clnt.Stats.NetIntf == nil {
		clnt.Stats.NetIntf = make(map[string]NetIntfInfo)
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 4 && (parts[2] == "inet" || parts[2] == "inet6") {
			ipv4 := parts[2] == "inet"
			intfname := parts[1]
			if info, ok := clnt.Stats.NetIntf[intfname]; ok {
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				clnt.Stats.NetIntf[intfname] = info
			} else {
				info := NetIntfInfo{}
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				clnt.Stats.NetIntf[intfname] = info
			}
		}
	}

	return
}

func (clnt *sshClient) getInterfaceInfo() (err error) {
	lines, err := clnt.RunCmd("/bin/cat /proc/net/dev")
	if err != nil {
		return
	}

	if clnt.Stats.NetIntf == nil {
		return
	} // should have been here already

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 17 {
			intf := strings.TrimSpace(parts[0])
			intf = strings.TrimSuffix(intf, ":")
			if info, ok := clnt.Stats.NetIntf[intf]; ok {
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
				clnt.Stats.NetIntf[intf] = info
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

func (clnt *sshClient) getCPU() (err error) {
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
	if clnt.Stats.preCPU.Total == 0 { // having no pre raw cpu data
		goto END
	}

	total = float32(nowCPU.Total - clnt.Stats.preCPU.Total)
	clnt.Stats.CPU.User = float32(nowCPU.User-clnt.Stats.preCPU.User) / total * 100
	clnt.Stats.CPU.Nice = float32(nowCPU.Nice-clnt.Stats.preCPU.Nice) / total * 100
	clnt.Stats.CPU.System = float32(nowCPU.System-clnt.Stats.preCPU.System) / total * 100
	clnt.Stats.CPU.Idle = float32(nowCPU.Idle-clnt.Stats.preCPU.Idle) / total * 100
	clnt.Stats.CPU.Iowait = float32(nowCPU.Iowait-clnt.Stats.preCPU.Iowait) / total * 100
	clnt.Stats.CPU.Irq = float32(nowCPU.Irq-clnt.Stats.preCPU.Irq) / total * 100
	clnt.Stats.CPU.SoftIrq = float32(nowCPU.SoftIrq-clnt.Stats.preCPU.SoftIrq) / total * 100
	clnt.Stats.CPU.Guest = float32(nowCPU.Guest-clnt.Stats.preCPU.Guest) / total * 100
END:
	clnt.Stats.preCPU = nowCPU
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
