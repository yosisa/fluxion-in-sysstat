package main

import (
	"os"
	"time"

	"github.com/yosisa/fluxion/message"
	"github.com/yosisa/fluxion/plugin"
)

var hostname, _ = os.Hostname()

type emitFunc func(string, map[string]interface{})

type Emitter interface {
	Emit(emitFunc) error
}

type EmitterFunc func(emitFunc) error

func (f EmitterFunc) Emit(emit emitFunc) error {
	return f(emit)
}

type Config struct {
	Tag          string   `toml:"tag"`
	Interval     string   `toml:"interval"`
	Processes    []string `toml:"processes"`
	Disks        []string `toml:"disks"`
	DiskInterval string   `toml:"disk_interval"`
}

type SysStatInput struct {
	env          *plugin.Env
	conf         *Config
	tagPrefix    string
	interval     time.Duration
	diskInterval time.Duration
	emitters     []Emitter
	closeCh      chan bool
}

func (p *SysStatInput) Init(env *plugin.Env) (err error) {
	p.env = env
	p.conf = &Config{}
	if err = env.ReadConfig(p.conf); err != nil {
		return
	}
	if p.conf.Tag != "" {
		p.tagPrefix = p.conf.Tag + "."
	}
	if p.interval, err = parseDuration(p.conf.Interval, 5*time.Second); err != nil {
		return
	}
	p.emitters = []Emitter{
		EmitterFunc(EmitMemory),
		EmitterFunc(EmitLoadAvg),
		&CPUStat{},
		&NetStat{},
	}
	if len(p.conf.Processes) > 0 {
		p.emitters = append(p.emitters, NewProcessStat(p.conf.Processes))
	}
	if len(p.conf.Disks) > 0 {
		if p.diskInterval, err = parseDuration(p.conf.DiskInterval, 1*time.Minute); err != nil {
			return
		}
	}
	return
}

func (p *SysStatInput) Start() error {
	go func() {
		var diskTick <-chan time.Time
		diskUsage := DiskUsage{Paths: p.conf.Disks}
		if p.diskInterval > 0 {
			diskTick = time.Tick(p.diskInterval)
			diskUsage.Emit(p.emit)
		}

		tick := time.Tick(p.interval)
		p.EmitStat()
		for {
			select {
			case <-tick:
				p.EmitStat()
			case <-diskTick:
				diskUsage.Emit(p.emit)
			}
		}
	}()
	return nil
}

func (p *SysStatInput) Close() error {
	return nil
}

func (p *SysStatInput) EmitStat() {
	for _, emitter := range p.emitters {
		err := emitter.Emit(p.emit)
		if err != nil {
			p.env.Log.Error(err)
		}
	}
}

func (p *SysStatInput) emit(tag string, v map[string]interface{}) {
	v["host"] = hostname
	p.env.Emit(message.NewEvent(p.tagPrefix+tag, v))
}

func main() {
	plugin.New("in-sysstat", func() plugin.Plugin {
		return &SysStatInput{}
	}).Run()
}

func parseDuration(s string, d time.Duration) (time.Duration, error) {
	if s == "" {
		return d, nil
	}
	return time.ParseDuration(s)
}
