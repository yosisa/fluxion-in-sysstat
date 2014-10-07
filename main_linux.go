package main

import (
	"strings"

	"github.com/shirou/gopsutil"
)

func EmitMemory(emit emitFunc) error {
	mem, err := gopsutil.VirtualMemory()
	if err != nil {
		return err
	}
	emit("memory", map[string]interface{}{
		"total":        mem.Total,
		"available":    mem.Available,
		"used":         mem.Used,
		"used_percent": mem.UsedPercent,
		"free":         mem.Free,
		"active":       mem.Active,
		"inactive":     mem.Inactive,
		"buffers":      mem.Buffers,
		"cached":       mem.Cached,
		"wired":        mem.Wired,
		"shared":       mem.Shared,
	})
	return nil
}

func EmitLoadAvg(emit emitFunc) error {
	load, err := gopsutil.LoadAvg()
	if err != nil {
		return err
	}
	emit("load", map[string]interface{}{
		"load1":  load.Load1,
		"load5":  load.Load5,
		"load15": load.Load15,
	})
	return nil
}

type ProcessStat struct {
	targets map[string]bool
}

func NewProcessStat(names []string) *ProcessStat {
	targets := make(map[string]bool)
	for _, name := range names {
		targets[name] = true
	}
	return &ProcessStat{targets: targets}
}

func (s *ProcessStat) Emit(emit emitFunc) error {
	pids, err := gopsutil.Pids()
	if err != nil {
		return err
	}

	for _, pid := range pids {
		p, err := gopsutil.NewProcess(pid)
		if err != nil {
			continue
		}
		cmd, err := p.Cmdline()
		if err != nil {
			continue
		}
		// Process.Name() doesn't contain full name. It maybe shortened.
		name := strings.SplitN(cmd, " ", 2)[0]
		if _, ok := s.targets[name]; !ok {
			continue
		}
		mem, err := p.MemoryInfoEx()
		if err != nil {
			continue
		}
		cpu, err := p.CPUTimes()
		if err != nil {
			continue
		}

		emit("process", map[string]interface{}{
			"pid":      p.Pid,
			"name":     name,
			"cmd":      cmd,
			"rss":      mem.RSS,
			"vms":      mem.VMS,
			"shared":   mem.Shared,
			"cpu_time": cpu.User + cpu.System,
		})
	}
	return nil
}
