package main

import (
	"time"

	"github.com/shirou/gopsutil"
)

type NetStat struct {
	old []gopsutil.NetIOCountersStat
	t   time.Time
}

func (s *NetStat) Emit(emit emitFunc) error {
	stats, err := gopsutil.NetIOCounters(true)
	if err != nil {
		return err
	}
	now := time.Now()

	if s.old == nil {
		// First time
		s.old = stats
		s.t = now
		return nil
	}

	sec := now.Sub(s.t).Seconds()
	persec := func(n uint64) uint64 {
		return uint64(float64(n) / sec)
	}

	for i, old := range s.old {
		cur := stats[i]
		if cur.Name != old.Name {
			// Find same interface if interface names are not match
			for _, v := range stats {
				if v.Name == old.Name {
					cur = v
					break
				}
			}
			// Skip if same interface is not found
			if cur.Name != old.Name {
				continue
			}
		}

		emit("net", map[string]interface{}{
			"name":       cur.Name,
			"rx_bytes":   persec(cur.BytesRecv - old.BytesRecv),
			"rx_packets": persec(cur.PacketsRecv - old.PacketsRecv),
			"rx_errors":  persec(cur.Errin - old.Errin),
			"rx_drop":    persec(cur.Dropin - old.Dropin),
			"tx_bytes":   persec(cur.BytesSent - old.BytesSent),
			"tx_packets": persec(cur.PacketsSent - old.PacketsSent),
			"tx_errors":  persec(cur.Errout - old.Errout),
			"tx_drop":    persec(cur.Dropout - old.Dropout),
		})
	}
	s.old = stats
	s.t = now
	return nil
}
