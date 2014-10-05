package main

import "github.com/shirou/gopsutil"

type DiskUsage struct {
	Paths []string
}

func (d *DiskUsage) Emit(emit emitFunc) error {
	for _, path := range d.Paths {
		usage, err := gopsutil.DiskUsage(path)
		if err != nil {
			continue
		}
		emit("disk_usage", map[string]interface{}{
			"path":               usage.Path,
			"total":              usage.Total,
			"free":               usage.Free,
			"used":               usage.Used,
			"used_percent":       usage.UsedPercent,
			"inode_total":        usage.InodesTotal,
			"inode_used":         usage.InodesUsed,
			"inode_free":         usage.InodesFree,
			"inode_used_percent": usage.InodesUsedPercent,
		})
	}
	return nil
}
