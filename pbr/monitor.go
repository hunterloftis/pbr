package pbr

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Monitor monitors several goroutines rendering stuff
type Monitor struct {
	Progress chan int
	Results  chan []float64

	cancel  chan interface{}
	active  int
	samples *safeCount
	nanos   int
	start   int64
}

type safeCount struct {
	sync.Mutex
	count int
}

// NewMonitor creates a new Monitor
func NewMonitor() *Monitor {
	return &Monitor{
		Progress: make(chan int, 32),
		Results:  make(chan []float64),
		cancel:   make(chan interface{}),
		samples:  &safeCount{},
		start:    time.Now().UnixNano(), // TODO: just use time.*, no need for ns conversion?
	}
}

// Samples returns the total number of samples
// TODO: reader lock?
func (m *Monitor) Samples() int {
	m.samples.Lock()
	defer m.samples.Unlock()
	return m.samples.count
}

// Nano returns the nanoseconds since Monitoring started
func (m *Monitor) Nano() int {
	return int(time.Now().UnixNano() - m.start)
}

// Active returns the number of active samplers/workers
func (m *Monitor) Active() int {
	return m.active
}

// Stop stops all workers
func (m *Monitor) Stop() {
	if m.Stopped() {
		return
	}
	close(m.cancel)
}

// SetInterrupt stops on interrupt or term signals
func (m *Monitor) SetInterrupt(callbacks ...func()) {
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt
		m.Stop()
		for _, f := range callbacks {
			f()
		}
	}()
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
			m.Progress <- total
			select {
			case <-m.cancel:
				m.active--
				m.Results <- s.Pixels()
				return
			default:
			}
		}
	}()
}
