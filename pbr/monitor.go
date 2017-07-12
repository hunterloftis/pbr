package pbr

import "runtime"

// Monitor monitors several goroutines rendering stuff
type Monitor struct {
	Sampler  *Sampler
	Renderer *Renderer
	samplers []*Sampler
}

// Start creates workers and starts monitoring
func (m *Monitor) Start(workers int) (update chan []float64, done chan []interface{}) {
	if workers == 0 {
		workers = runtime.NumCPU()
	}
	update = make(chan []float64)
	done = make(chan []interface{})
	m.samplers = []*Sampler{}
	for len(m.samplers) < workers {
		m.samplers = append(m.samplers, m.Sampler.Clone())
	}

	return
}
