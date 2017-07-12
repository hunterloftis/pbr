package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	var target, focus *pbr.Vector3
	position := pbr.Vector3{0, 0, 0}
	sky := pbr.Vector3{40, 50, 60}
	in := os.Args[1]
	out := flag.String("out", "render.png", "Output png filename")
	heat := flag.String("heat", "", "Heatmap png filename")
	workers := flag.Int("workers", runtime.NumCPU(), "Concurrency level")
	samples := flag.Float64("samples", math.Inf(1), "Max samples per pixel")
	adapt := flag.Int("adapt", 4, "Adaptive sampling; 0=off, 3=medium, 5=high")
	bounces := flag.Int("bounces", 10, "Maximum light bounces")
	profile := flag.Bool("profile", false, "Record performance into profile.pprof")
	pano := flag.String("pano", "", "Panoramic environment map hdr (radiosity) file")
	lens := flag.Float64("lens", 50, "Camera focal length in mm")
	exposure := flag.Float64("exposure", 1, "Exposure multiplier")
	fStop := flag.Float64("fstop", 4, "Camera f-stop")
	flag.Var(&position, "position", "Camera position")
	flag.Var(target, "target", "Camera target location")
	flag.Var(focus, "focus", "Camera focus location")
	flag.Var(&sky, "sky", "Ambient sky lighting")
	flag.Parse()

	if *profile {
		f, _ := os.Create("profile.pprof")
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	xml, _ := ioutil.ReadFile(in)
	scene := pbr.ColladaScene(xml)
	camera := pbr.NewCamera(1280, 720, pbr.CameraConfig{
		Lens:     (*lens) / 1000.0,
		Position: position,
		Target:   target,
		Focus:    focus,
		FStop:    *fStop,
	})
	renderer := pbr.CamRenderer(camera, pbr.RenderConfig{
		Exposure: *exposure,
	})

	scene.SetSky(sky, pbr.Vector3{})
	if len(*pano) > 0 {
		hdr, _ := os.Open(*pano)
		defer hdr.Close()
		scene.SetPano(hdr, 100) // TODO: read radiosity info or allow it from the command line
	}

	total := float64(renderer.Width*renderer.Height) * (*samples)
	interrupt := make(chan os.Signal, 2)
	update := make(chan float64)
	done := make(chan []interface{})
	results := make(chan []float64)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	for i := 0; i < *workers; i++ {
		go func(i int) {
			fmt.Println("creating worker:", i)
			s := pbr.NewSampler(camera, scene, pbr.SamplerConfig{
				Bounces: *bounces,
				Adapt:   *adapt,
			})
			for i := 0.0; i < total; i += float64(s.SampleFrame()) {
				select {
				case <-done:
					fmt.Println("worker got done signal:", i)
					break
				default:
					update <- i
				}
			}
			results <- s.Pixels()
		}(i)
	}

	pending := *workers
workerLoop:
	for pending > 0 {
		select {
		case s := <-update:
			fmt.Printf("%v samples\n", s)
		case r := <-results:
			fmt.Println("merging...")
			renderer.Merge(r)
			pending--
		case <-interrupt:
			fmt.Println("interrupting")
			break workerLoop
		case <-done:
			break workerLoop
		default:
		}
	}
	fmt.Println("closing done")
	close(done)

	writePNG(*out, renderer.Rgb())
	fmt.Printf("-> %v\n", *out)
	if len(*heat) > 0 {
		writePNG(*heat, renderer.Heat())
		fmt.Printf("-> %v\n", *heat)
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
