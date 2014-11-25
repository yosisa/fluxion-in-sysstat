package main

import "github.com/shirou/gopsutil"

type CPUStat struct {
	old []gopsutil.CPUTimesStat
}

func (s *CPUStat) Emit(emit emitFunc) error {
	stats, err := gopsutil.CPUTimes(false)
	if err != nil {
		return err
	}

	if s.old == nil {
		// First time
		s.old = stats
		return nil
	}

	for i, old := range s.old {
		cur := stats[i]
		user := cur.User - old.User
		nice := cur.Nice - old.Nice
		system := cur.System - old.System
		idle := cur.Idle - old.Idle
		iowait := cur.Iowait - old.Iowait
		irq := cur.Irq - old.Irq
		softirq := cur.Softirq - old.Softirq
		stolen := cur.Stolen - old.Stolen

		total := user + nice + system + idle + iowait + irq + softirq + stolen
		idle_percent := idle / total * 100
		emit("cpu", map[string]interface{}{
			"cpu":     old.CPU,
			"used":    100 - idle_percent,
			"user":    user / total * 100,
			"nice":    nice / total * 100,
			"system":  system / total * 100,
			"idle":    idle_percent,
			"iowait":  iowait / total * 100,
			"irq":     irq / total * 100,
			"softirq": softirq / total * 100,
			"stolen":  stolen / total * 100,
		})
	}
	s.old = stats
	return nil
}
