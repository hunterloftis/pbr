package pbr

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"
)

// Cli is an abstraction for executing a render via a terminal.
type Cli struct {
	scene    *Scene
	cam      *Camera
	renderer *Renderer
}

type statistic struct {
	sync.Mutex
	samples int
	start   int64
}

// CliRunner constructs a CLI from pointers to a scene, camera, and renderer.
func CliRunner(scene *Scene, cam *Camera, renderer *Renderer) Cli {
	c := Cli{scene, cam, renderer}
	return c
}

// Render parses command-line flags and creates
// workers to render its given scene, in parallel, from the point-of-view of its given camera.
// Unless given a -samples argument, it renders increasingly high-fidelity images
// until it's interrupted by a signal (like SIGINT / Ctrl+C).
// Once it receives a signal or passes its sampling threshold, it writes results to PNGs.
func (c Cli) Render() {
	out := flag.String("out", "render.png", "Output png filename")
	heat := flag.String("heat", "", "Heatmap png filename")
	workers := flag.Int("workers", runtime.NumCPU(), "Concurrency level")
	samples := flag.Float64("samples", math.Inf(1), "Max samples per pixel")
	profile := flag.Bool("profile", false, "Record performance into profile.pprof")
	adapt := flag.Int("adapt", 4, "Adaptive sampling; 0=off, 3=medium, 5=high")
	bounces := flag.Int("bounces", 10, "Maximum light bounces")
	flag.Parse()

	if *profile {
		f, _ := os.Create("profile.pprof")
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	working := make(chan struct{})
	interrupted := make(chan os.Signal, 2)
	results := make(chan []float64)

	signal.Notify(interrupted, os.Interrupt, syscall.SIGTERM)
	go func() { <-interrupted; fmt.Printf("\n => Interrupting...\n"); close(working) }()

	fmt.Printf("Rendering (%v workers, %v bounces, per-pixel samples=%v, adapt=%v)\n", *workers, *bounces, *samples, *adapt)
	stat := statistic{start: time.Now().UnixNano()}
	for i := 0; i < *workers; i++ {
		go c.worker(&stat, *samples, *bounces, *adapt, working, results)
	}
	for i := 0; i < *workers; i++ {
		c.renderer.Merge(<-results)
	}
	fmt.Printf(" => Writing...\n")
	writePNG(*out, c.renderer.Rgb())
	fmt.Printf("    RGB: %v\n", *out)
	if len(*heat) > 0 {
		writePNG(*heat, c.renderer.Heat())
		fmt.Printf("    Heat: %v\n", *heat)
	}
}

// TODO: this is a really long argument list
func (c Cli) worker(stat *statistic, max float64, bounces int, adapt int, done <-chan struct{}, results chan<- []float64) {
	sampler := NewSampler(c.cam, c.scene, bounces, adapt)
	pixels := sampler.Width * sampler.Height
	for {
		select {
		case <-done:
			results <- sampler.Pixels()
			return
		default:
			if float64(stat.samples/pixels) >= max {
				fmt.Printf(" => sample limit\n")
				results <- sampler.Pixels()
				return
			}
			samples := sampler.SampleFrame()
			stat.Lock()
			stat.samples += samples
			stat.Unlock()
			ms := float64(time.Now().UnixNano()-stat.start) * 1e-6
			fmt.Printf(" => %v samples (%v / pixel, %v / ms)\n", stat.samples, stat.samples/pixels, int(float64(stat.samples)/ms))
		}
	}
}

func writePNG(file string, i image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, i)
}
