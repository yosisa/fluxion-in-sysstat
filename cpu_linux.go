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
		user := old.User - stats[i].User
		nice := old.Nice - stats[i].Nice
		system := old.System - stats[i].System
		idle := old.Idle - stats[i].Idle
		iowait := old.Iowait - stats[i].Iowait
		irq := old.Irq - stats[i].Irq
		softirq := old.Softirq - stats[i].Softirq
		stolen := old.Stolen - stats[i].Stolen

		total := user + nice + system + idle + iowait + irq + softirq + stolen
		emit("cpu", map[string]interface{}{
			"cpu":     old.CPU,
			"user":    user / total * 100,
			"nice":    nice / total * 100,
			"system":  system / total * 100,
			"idle":    idle / total * 100,
			"iowait":  iowait / total * 100,
			"irq":     irq / total * 100,
			"softirq": softirq / total * 100,
			"stolen":  stolen / total * 100,
		})
	}
	s.old = stats
	return nil
}
