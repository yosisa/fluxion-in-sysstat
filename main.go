package main

import (
	"time"

	"github.com/yosisa/fluxion/event"
	"github.com/yosisa/fluxion/plugin"
)

type Config struct {
	Tag      string `codec:"tag"`
	Interval string `codec:"interval"`
}

type SysStatInput struct {
	env       *plugin.Env
	conf      *Config
	tagPrefix string
	interval  time.Duration
	closeCh   chan bool
}

func (p *SysStatInput) Name() string {
	return "in-sysstat"
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
	if p.conf.Interval == "" {
		p.interval = 6 * time.Second
	} else {
		p.interval, err = time.ParseDuration(p.conf.Interval)
	}
	return
}

func (p *SysStatInput) Start() error {
	go func() {
		tick := time.Tick(p.interval)
		for {
			p.EmitStat()
			<-tick
		}
	}()
	return nil
}

func (p *SysStatInput) EmitStat() {
	v, err := GetMemory()
	if err == nil {
		p.env.Emit(event.NewRecord(p.tagPrefix+"memory", v))
	}
}

func main() {
	plugin.New(func() plugin.Plugin {
		return &SysStatInput{}
	}).Run()
}
