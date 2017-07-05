package pbr

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

// Cli is a command line abstraction for rendering
type Cli struct {
	scene    *Scene
	cam      *Camera
	renderer *Renderer
}

type statistic struct {
	sync.Mutex
	samples int
}

// NewCLI makes a new CLI
func NewCLI(scene *Scene, cam *Camera, renderer *Renderer) Cli {
	c := Cli{scene, cam, renderer}
	return c
}

// Start starts rendering based on CLI parameters
func (c Cli) Start() {
	out := flag.String("out", "render.png", "Output png filename")
	heat := flag.String("heat", "", "Heatmap png filename")
	workers := flag.Int("workers", runtime.NumCPU(), "Concurrency level")
	samples := flag.Float64("samples", math.Inf(1), "Max samples per pixel")
	flag.Parse()

	working := make(chan struct{})
	interrupted := make(chan os.Signal, 2)
	results := make(chan []float64)

	signal.Notify(interrupted, os.Interrupt, syscall.SIGTERM)
	go func() { <-interrupted; fmt.Println(""); close(working) }()

	fmt.Printf("Rendering (%v workers, %v samples/pixel)\n", *workers, *samples)
	stat := statistic{}
	for i := 0; i < *workers; i++ {
		go c.worker(&stat, *samples, working, results)
	}
	for i := 0; i < *workers; i++ {
		c.renderer.Merge(<-results)
	}
	fmt.Printf(" => Writing...\n")
	c.renderer.WriteRGB(*out)
	fmt.Printf("    RGB: %v\n", *out)
	if len(*heat) > 0 {
		c.renderer.WriteHeat(*heat)
		fmt.Printf("    Heat: %v\n", *heat)
	}
}

func (c Cli) worker(stat *statistic, max float64, done <-chan struct{}, results chan<- []float64) {
	sampler := NewSampler(c.cam, c.scene, 10, 3)
	pixels := sampler.Width * sampler.Height
	for {
		samples := sampler.SampleFrame()
		stat.Lock()
		stat.samples += samples
		stat.Unlock()
		fmt.Printf(" => %v samples (%v / pixel)\n", stat.samples, stat.samples/pixels)
		select {
		case <-done:
			results <- sampler.Pixels()
			return
		default:
			if float64(stat.samples/pixels) >= max {
				results <- sampler.Pixels()
				return
			}
		}
	}
}
