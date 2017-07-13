package pbr

import "sync"

// Monitor monitors several goroutines rendering stuff
type Monitor struct {
	Progress chan int
	Results  chan []float64

	cancel  chan interface{}
	active  int
	samples *safeCount
}

type safeCount struct {
	sync.Mutex
	count int
}

// NewMonitor creates a new Monitor
func NewMonitor() *Monitor {
	return &Monitor{
		Progress: make(chan int, 10),
		Results:  make(chan []float64),
		cancel:   make(chan interface{}),
		samples:  &safeCount{},
	}
}

// Active returns the number of active samplers/workers
func (m *Monitor) Active() int {
	return m.active
}

// Stop stops all workers
func (m *Monitor) Stop() {
	select {
	case <-m.cancel:
		return
	default:
		close(m.cancel)
	}
}

// Stopped returns whether or not this is stopped
func (m *Monitor) Stopped() bool {
	select {
	case <-m.cancel:
		return true
	default:
		return false
	}
}

// AddSampler creates a new worker with that sampler
func (m *Monitor) AddSampler(s *Sampler) {
	m.active++
	go func() {
		for {
			frame := s.SampleFrame()
			m.samples.Lock()
			m.samples.count += frame
			total := m.samples.count // TODO: do I need to do this or can I safely read after unlocking?
			m.samples.Unlock()
			select {
			case <-m.cancel:
				m.Results <- s.Pixels()
				m.active--
				return
			case m.Progress <- total:
			default:
			}
		}
	}()
}
