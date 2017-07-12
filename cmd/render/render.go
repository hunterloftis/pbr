package main

import (
	"flag"
	"image"
	"image/png"
	"os"
	"runtime/pprof"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	// in := os.Args[1]
	// out := flag.String("out", "render.png", "Output png filename")
	// heat := flag.String("heat", "", "Heatmap png filename")
	// workers := flag.Int("workers", runtime.NumCPU(), "Concurrency level")
	// samples := flag.Float64("samples", math.Inf(1), "Max samples per pixel")
	// adapt := flag.Int("adapt", 4, "Adaptive sampling; 0=off, 3=medium, 5=high")
	// bounces := flag.Int("bounces", 10, "Maximum light bounces")
	profile := flag.Bool("profile", false, "Record performance into profile.pprof")
	// sky := flag.String("sky", "40,50,60", "Ambient sky lighting RGB") // TODO: setter/getter
	// pano := flag.String("pano", "", "Panoramic environment map hdr (radiosity) file")
	// cam := flag.String("cam", "1,1,1", "Camera position")            // TODO: setter/getter
	// target := flag.String("look", "0,0,0", "Camera target location") // TODO: setter/getter
	// focus := flag.String("focus", "0,0,0", "Camera focus location")  // TODO: setter/getter
	// lens := flag.Float64("lens", 50, "Camera focal length in mm")    // TODO: setter/getter to convert to 0.050
	// exposure := flag.Float64("exposure", 1, "Exposure multiplier")
	// fStop := flag.Float64("fstop", 4, "Camera f-stop")
	position := pbr.Vector3{1, 1, 1}
	flag.Var(&position, "position", "Camera position")
	flag.Parse()

	if *profile {
		f, _ := os.Create("profile.pprof")
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// xml, _ := ioutil.ReadFile(in)
	// scene := pbr.ColladaScene(xml)
	// _ = sky

	// // scene.SetSky(sky)
	// if len(*pano) > 0 {
	// 	hdr, _ := os.Open(*pano)
	// 	defer hdr.Close()
	// 	scene.SetPano(hdr, 100) // TODO: read radiosity info or allow it from the command line
	// }
	// camera := pbr.NewCamera(1280, 720, pbr.CameraConfig{
	// 	Lens:     *lens,
	// 	Position: *cam,
	// 	Target:   *target,
	// 	Focus:    *focus,
	// 	FStop:    *fStop,
	// })
	// sampler := pbr.NewSampler(camera, scene, pbr.SamplerConfig{
	// 	Bounces: *bounces,
	// 	Samples: *samples,
	// 	Adapt:   *adapt,
	// })
	// renderer := pbr.CamRenderer(camera, pbr.RenderConfig{
	// 	Exposure: 1,
	// })
	// monitor := pbr.NewMonitor(sampler, renderer, pbr.MonitorConfig{
	// 	Workers: *workers,
	// })

	// monitor.Start()

	// for pp := range <-monitor.C {
	// 	fmt.Printf("%v samples / pixel\n", pp)
	// }
	// writePNG(*out, renderer.Rgb())
	// fmt.Printf("-> %v\n", *out)
	// if len(*heat) > 0 {
	// 	writePNG(*heat, renderer.Heat())
	// 	fmt.Printf("-> %v\n", *heat)
	// }
}

func writePNG(file string, i image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, i)
}
