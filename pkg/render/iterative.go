package render

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Limits struct {
	DumpPeriod time.Duration
	Time       time.Duration
	Frames     int
}

func NewLimits(d int, t int, f int) *Limits {
	if d == 0 {
		d = math.MaxInt32
	}
	if t == 0 {
		t = math.MaxInt32
	}
	return &Limits{
		DumpPeriod: time.Duration(d) * time.Second,
		Time:       time.Duration(t) * time.Second,
		Frames:     f,
	}
}

func Iterative(scene *Scene, limits *Limits, file string, width, height, depth int, direct bool) error {
	kill := make(chan os.Signal, 2)
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	if limits == nil {
		limits = NewLimits(6, 0, 0)
	}

	frame := scene.Render(width, height, depth, direct)
	defer frame.Stop()
	ticker := time.NewTicker(limits.DumpPeriod)
	defer ticker.Stop()
	limiter := time.NewTicker(limits.Time)
	defer limiter.Stop()
	framelim := time.NewTicker(1 * time.Second)
	if limits.Frames == 0 {
		framelim.Stop()
	} else {
		defer framelim.Stop()
	}

	start := time.Now().UnixNano()
	max := 0
	fmt.Printf("\nRendering %v (Ctrl+C to end)", file)

	for frame.Active() {
		select {
		case <-kill:
			frame.Stop()
		case <-limiter.C:
			frame.Stop()
			fmt.Printf("\nTime limit reached.\n")
		case <-framelim.C:
			if _, n := frame.Sample(); n > limits.Frames {
				frame.Stop()
				fmt.Print("\nFrame limit reached.\n")
			}
		case <-ticker.C:
			if sample, n := frame.Sample(); n > max {
				max = n
				fmt.Print(".")
				if err := writePng(file, sample.Image()); err != nil {
					return err
				}
			}
		}
	}

	stop := time.Now().UnixNano()
	sample, frames := frame.Sample()
	total := sample.Total()
	p := message.NewPrinter(language.English)
	secs := float64(stop-start) / 1e9
	sps := math.Round(float64(total) / secs) // TODO: rename to pixels/sec for clarity
	p.Printf("\n%v samples in %.1f seconds (%.0f samples/sec)\n", total, secs, sps)
	p.Printf("\n%v frames (%.1f frames/sec)\n", frames, float64(frames)/secs)
	if err := writePng(file, sample.Image()); err != nil {
		return err
	}

	return nil
}

func writePng(filename string, im image.Image) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, im)
}
