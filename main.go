package main

import (
	"time"

	"github.com/yosisa/fluxion/event"
	"github.com/yosisa/fluxion/plugin"
)

type emitFunc func(string, map[string]interface{})

type Emitter interface {
	Emit(emitFunc) error
}

type EmitterFunc func(emitFunc) error

func (f EmitterFunc) Emit(emit emitFunc) error {
	return f(emit)
}

type Config struct {
	Tag       string   `codec:"tag"`
	Interval  string   `codec:"interval"`
	Processes []string `codec:"processes"`
}

type SysStatInput struct {
	env       *plugin.Env
	conf      *Config
	tagPrefix string
	interval  time.Duration
	emitters  []Emitter
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
	p.emitters = []Emitter{
		EmitterFunc(EmitMemory),
	}
	if len(p.conf.Processes) > 0 {
		p.emitters = append(p.emitters, NewProcessStat(p.conf.Processes))
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
	for _, emitter := range p.emitters {
		err := emitter.Emit(func(tag string, v map[string]interface{}) {
			p.env.Emit(event.NewRecord(p.tagPrefix+tag, v))
		})
		if err != nil {
			p.env.Log.Error(err)
		}
	}
}

func main() {
	plugin.New(func() plugin.Plugin {
		return &SysStatInput{}
	}).Run()
}
