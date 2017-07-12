package pbr

import "runtime"

// Monitor monitors several goroutines rendering stuff
type Monitor struct {
	C <-chan []float64
	MonitorConfig
	samplers []*Sampler
}

// MonitorConfig configures a Monitor
type MonitorConfig struct {
	Workers int
}

// NewMonitor creates a new monitor
func NewMonitor(sampler *Sampler, renderer *Renderer, config ...MonitorConfig) *Monitor {
	c := config[0]
	if c.Workers == 0 {
		c.Workers = runtime.NumCPU() // use one worker per CPU by default
	}
	m := Monitor{
		C:             make(<-chan []float64),
		MonitorConfig: c,
		samplers:      []*Sampler{sampler},
	}
	for len(m.samplers) < m.Workers {
		m.samplers = append(m.samplers, sampler.Clone())
	}
	return &m
}

// Start creates workers and starts monitoring
func (m *Monitor) Start() {

}
