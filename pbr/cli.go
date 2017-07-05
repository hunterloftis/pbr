package pbr

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// Cli is a command line abstraction for rendering
type Cli struct {
	scene    *Scene
	cam      *Camera
	renderer *Renderer
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
	samples := flag.Float64("samples", math.Inf(1), "Average per pixel samples to take.")
	heat := flag.String("heat", "", "Heatmap png filename.")
	flag.Parse()

	working := make(chan struct{})
	interrupted := make(chan os.Signal, 2)
	results := make(chan []float64)

	signal.Notify(interrupted, os.Interrupt, syscall.SIGTERM)
	go func() { <-interrupted; close(working) }()

	fmt.Printf("Rendering (concurrency: %v, per pixel sample cutoff: %v)\n", *concurrency, *samples)
	for i := 0; i < *concurrency; i++ {
		go c.worker(i, working, results) // instantiate a worker ala https://play.golang.org/p/Sfx1JL_6K2
	}
	<-working

	fmt.Printf("\nMerging...\n")
	for i := 0; i < *concurrency; i++ {
		fmt.Printf("Merging %v\n", i)
		c.renderer.Merge(<-results)
	}
	fmt.Printf("Writing...\n")
	c.renderer.WriteRGB(*out)
	if len(*heat) > 0 {
		c.renderer.WriteHeat(*heat)
	}
	fmt.Printf("\n -> %v\n", *out)
}

func (c Cli) worker(id int, done <-chan struct{}, results chan<- []float64) {
	// TODO: implement Sampler.Clone() and Clone() a passed-in Sampler instead of instantiating one here?
	sampler := NewSampler(c.cam, c.scene, 10, 3)
	for {
		fmt.Printf("Worker %v...\n", id)
		sampler.SampleFrame()
		// TODO: periodically merge & write?
		select {
		case <-done:
			fmt.Printf("Worker %v exiting...\n", id)
			results <- sampler.Pixels()
			return
		default:
		}
	}
}
