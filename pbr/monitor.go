package pbr

import (
	"fmt"
	"sync"
	"time"
)

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
		Progress: make(chan int),
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
	if m.Stopped() {
		return
	}
	close(m.cancel)
}

// Stopped returns whether or not this is stopped
func (m *Monitor) Stopped() bool {
	select {
	case <-m.cancel:
		fmt.Println("Stopped - true")
		return true
	default:
		fmt.Println("Stopped - false")
		return false
	}
}

// AddSampler creates a new worker with that sampler
func (m *Monitor) AddSampler(s *Sampler) {
	m.active++
	go func() {
		for {
			fmt.Println("Start SampleFrame()")
			start := time.Now().UnixNano()
			frame := s.SampleFrame()
			secs := float64(time.Now().UnixNano()-start) * 1e-9
			fmt.Println("End SampleFrame(), seconds taken:", secs)
			m.samples.Lock()
			m.samples.count += frame
			total := m.samples.count // TODO: do I need to do this or can I safely read after unlocking?
			m.samples.Unlock()
			fmt.Println("<- send progress")
			m.Progress <- total
			fmt.Println("progress sent")
			select {
			case <-m.cancel:
				fmt.Println("<- Send pixels to m.Results")
				m.Results <- s.Pixels()
				fmt.Println("pixels sent")
				m.active--
				return
			default:
				fmt.Println("Restart sampling")
			}
		}
	}()
}
