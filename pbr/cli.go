package pbr

import (
	"flag"
	"fmt"
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
	out := flag.String("out", "render.png", "Output png filename.")
	concurrency := flag.Int("frames", runtime.NumCPU(), "Number of frames to combine.")
	heat := flag.String("heat", "", "Heatmap png filename.")
	flag.Parse()

	working := make(chan struct{})
	interrupted := make(chan os.Signal, 2)
	results := make(chan []float64)

	signal.Notify(interrupted, os.Interrupt, syscall.SIGTERM)
	go func() { <-interrupted; close(working) }()

	fmt.Printf("Rendering (%v workers)\n", *concurrency)
	stat := statistic{}
	for i := 0; i < *concurrency; i++ {
		go c.worker(&stat, working, results) // instantiate a worker ala https://play.golang.org/p/Sfx1JL_6K2
	}
	<-working

	fmt.Printf("\n => Merging...\n", stat.samples)
	for i := 0; i < *concurrency; i++ {
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

func (c Cli) worker(stat *statistic, done <-chan struct{}, results chan<- []float64) {
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
		}
	}
}
