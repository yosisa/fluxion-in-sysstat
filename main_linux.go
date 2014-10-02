package main

import "github.com/shirou/gopsutil"

func GetMemory() (map[string]interface{}, error) {
	mem, err := gopsutil.VirtualMemory()
	if err != nil {
		return nil, err
	}
	v := map[string]interface{}{
		"total":       mem.Total,
		"available":   mem.Available,
		"used":        mem.Used,
		"usedPercent": mem.UsedPercent,
		"free":        mem.Free,
		"active":      mem.Active,
		"inactive":    mem.Inactive,
		"buffers":     mem.Buffers,
		"cached":      mem.Cached,
		"wired":       mem.Wired,
		"shared":      mem.Shared,
	}
	return v, nil
}
