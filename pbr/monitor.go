package pbr

// Monitor monitors several goroutines rendering stuff
type Monitor struct {
	Progress chan float64
	Results  chan []float64

	cancel chan interface{}
	active int
	count  int
}

// NewMonitor creates a new Monitor
func NewMonitor() *Monitor {
	return &Monitor{
		Progress: make(chan float64),
		Results:  make(chan []float64),
		cancel:   make(chan interface{}),
	}
}

// Active returns the number of active samplers/workers
func (m *Monitor) Active() int {
	return m.active
}

// Stop stops all workers
func (m *Monitor) Stop() {
	close(m.cancel)
}

// AddSampler creates a new worker with that sampler
func (m *Monitor) AddSampler(s *Sampler) {
	m.active++
	m.count++
	i := m.count
	go func() {
		for {
			s.SampleFrame()
			select {
			case <-m.cancel:
				m.Results <- s.Pixels()
				m.active--
				return
			default:
				m.Progress <- float64(i)
			}
		}
	}()
}
