package pbr

import (
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

// Render renders the samples in a Sampler to an Image.
// TODO: implement a de-noising filter; https://www.youtube.com/watch?v=Ee51bkOlbMw
type Render struct {
	width   int
	height  int
	adapt   float64
	bounces int
	direct  int
	branch  int

	camera *Camera
	scene  *Scene
	rnd    *rand.Rand

	// state
	active       bool
	count        uint
	cursor       uint
	meanVariance float64
	light        *rgb.Lightmap
	rays         uint
}

type sample struct {
	index  uint
	energy rgb.Energy
}

// NewRender creates a renderer referencing a Sampler.
func NewRender(s *Scene, c *Camera) *Render {
	return &Render{
		width:        c.Width(),
		height:       c.Height(),
		bounces:      8,
		adapt:        10,
		branch:       64,
		direct:       1,
		camera:       c,
		scene:        s,
		meanVariance: math.MaxFloat64,
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
		light:        rgb.Map(uint(c.Width()), uint(c.Height())),
	}
}

func (r *Render) SetBounces(b int) {
	r.bounces = b
}

func (r *Render) SetAdapt(a float64) {
	r.adapt = a
}

func (r *Render) SetDirect(d int) {
	r.direct = d
}

func (r *Render) SetBranch(b int) {
	r.branch = b
}

func (r *Render) Start() {
	r.active = true
	r.scene.Prepare()
	n := runtime.NumCPU()
	buffer := n * 4
	result := make(chan sample, buffer)             // TODO: use a chan []sample instead with a configurable number of samples to batch
	pixel := make(chan uint, int(r.adapt+1)*buffer) // TODO: use a chan []uint instead with a configurable number of pixel indices to batch
	planned := uint(0)
	for i := 0; i < n; i++ {
		r.sample(pixel, result)
		planned += r.next(pixel)
	}
	go func() {
		for {
			res := <-result
			r.light.Integrate(res.index, res.energy)
			r.count++
			for planned < r.count+uint(buffer) {
				planned += r.next(pixel)
			}
			if !r.active {
				close(pixel)
				return
			}
		}
	}()
}

func (r *Render) Stop() {
	r.active = false
}

func (r *Render) Active() bool {
	return r.active
}

func (r *Render) Count() uint {
	return r.count
}

func (r *Render) Size() uint {
	return uint(r.camera.Width() * r.camera.Height())
}

func (r *Render) Light() *rgb.Lightmap {
	return r.light
}

func (r *Render) Image(expose float64) image.Image {
	return r.light.Image(expose)
}

func (r *Render) Rays() uint {
	return r.rays
}

// TODO: sample in chunks of N samples, then send several results back at a time.
// Communication overhead is high.
// Maybe by row?
func (r *Render) sample(in <-chan uint, out chan<- sample) {
	size := uint(r.width * r.height)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	go func() {
		for {
			if p, ok := <-in; ok {
				i := p % size
				x, y := r.pixelAt(i)
				ray := r.camera.ray(x, y, rnd)
				e := r.trace2(ray, 0, rgb.Energy{1, 1, 1}, rnd)
				// e := r.trace(x, y, rnd)
				out <- sample{i, e}
			} else {
				return
			}
		}
	}()
}

func (r *Render) pixelAt(i uint) (x, y float64) {
	return float64(i % uint(r.width)), float64(i / uint(r.width))
}

// http://blog.yiningkarlli.com/2013/04/working-towards-importance-sampled-direct-lighting.html
// https://www.fxguide.com/featured/the-art-of-rendering/
// http://web.cs.wpi.edu/~emmanuel/courses/cs563/write_ups/zackw/realistic_raytracing.html
// https://www.cs.utexas.edu/~jthywiss/cs384g_proj4.shtml
// http://www.cs.cornell.edu/courses/cs6620/2009sp/Lectures/Lec10_MonteCarlo_web.pdf
// https://blender.stackexchange.com/questions/3256/what-is-branched-path-tracing-and-how-is-it-useful
// https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=1&cad=rja&uact=8&ved=0ahUKEwiX4MqS5dzYAhVlUd8KHRdJA6YQFggpMAA&url=https%3A%2F%2Fwww.cs.northwestern.edu%2F~jet%2FTeach%2F2003_1winAdvGraphics%2Fpresentations%2Fweek5.globillum.McCrory.ppt&usg=AOvVaw1FKrkeV713-6cozb9xFXeJ
// TODO: use https://blog.carlmjohnson.net/post/2016-11-27-how-to-use-go-generate/
// to flag in correctness assertions for testing (like dir.Len() == 1)
// http://www.cs.cornell.edu/courses/cs6630/2012sp/notes/07pathtr-notes.pdf
// http://blog.yiningkarlli.com/2013/04/importance-sampled-direct-lighting.html
// https://graphics.stanford.edu/courses/cs348b-03/lectures/mc-3.pdf
// TODO: take several samples at each bounce (since the ray intersections are the most expensive part)
// Maybe (maxdepth - depth)^2 or something, to get more accurate readings in the earlier bounces?
// Or maybe: (1 - gloss) * (maxdepth - depth) * (maxdepth - depth)
// Or maybe: 1 + (1 - gloss) * ((maxdepth - depth) / maxdepth) * r.indirect
// Or maybe: 1 + (1 - gloss) * r.indirect ... and then let russian roulette take care of early termination
// https://blender.stackexchange.com/questions/3256/what-is-branched-path-tracing-and-how-is-it-useful
func (r *Render) trace(x, y float64, rnd *rand.Rand) rgb.Energy {
	ray := r.camera.ray(x, y, rnd)
	signal := rgb.Energy{1, 1, 1}
	energy := rgb.Energy{0, 0, 0}

	for bounce := 0; bounce < r.bounces; bounce++ {
		hit := r.scene.Intersect(ray)
		r.rays++ // TODO: store count locally and merge later (since can't trust writing from multiple goroutines)
		if !hit.Ok {
			energy = energy.Merged(r.scene.EnvAt(ray), signal)
			break
		}
		point := ray.Moved(hit.Dist)
		normal, mat := hit.Surface.At(point)
		energy = energy.Merged(mat.Emit(), signal)
		dir, strength, diffused := mat.Bsdf(normal, ray.Dir, hit.Dist, rnd)
		if diffused {
			direct, coverage := r.traceDirect(point, normal, rnd)
			energy = energy.Merged(direct.Strength(mat.Color()), signal)
			signal = signal.Amplified(1 - coverage)
		}
		if signal = signal.Strength(strength).RandomGain(rnd); signal.Zero() { // "Russian Roulette"
			break
		}
		ray = geom.NewRay(point, dir)
	}
	return energy
}

func (r *Render) trace2(ray *geom.Ray3, depth int, signal rgb.Energy, rnd *rand.Rand) (energy rgb.Energy) {
	if depth >= r.bounces {
		return energy
	}
	hit := r.scene.Intersect(ray)
	r.rays++
	if !hit.Ok {
		energy = energy.Merged(r.scene.EnvAt(ray), signal)
		return energy
	}
	point := ray.Moved(hit.Dist)
	normal, mat := hit.Surface.At(point)
	energy = energy.Merged(mat.Emit(), signal)
	if signal = signal.RandomGain(rnd); signal.Zero() {
		return energy
	}
	branch := 1
	if depth == 0 {
		branch += int(float64(r.branch) * (mat.Roughness() + 0.25) / 1.25)
	}
	sum := rgb.Energy{}
	for i := 0; i < branch; i++ { // TODO: locally adapt to noise? Could skip the whole middle-man of noise tracking.
		dir, strength, diffused := mat.Bsdf(normal, ray.Dir, hit.Dist, rnd)
		sig := signal
		if diffused {
			direct, coverage := r.traceDirect(point, normal, rnd)
			sum = sum.Merged(direct.Strength(mat.Color()), sig)
			sig = sig.Amplified(1 - coverage)
		}
		next := geom.NewRay(point, dir)
		sum = sum.Plus(r.trace2(next, depth+1, sig.Strength(strength), rnd))
	}
	average := sum.Amplified(1 / float64(branch))
	return energy.Plus(average)
}

func (r *Render) traceDirect(point geom.Vector3, normal geom.Direction, rnd *rand.Rand) (energy rgb.Energy, coverage float64) {
	for i := 0; i < r.direct; i++ {
		light := r.scene.Light(rnd)
		shadow, solidAngle := light.Box().ShadowRay(point, rnd)
		cos := shadow.Dir.Cos(normal)
		if cos <= 0 {
			break
		}
		coverage += solidAngle
		hit := r.scene.Intersect(shadow)
		r.rays++
		if !hit.Ok {
			break
		}
		e := hit.Surface.Material().Emit().Amplified(solidAngle * cos / math.Pi)
		energy = energy.Plus(e)
	}
	return energy, coverage
}

func (r *Render) next(pixels chan<- uint) uint {
	count := uint(1)
	if r.adapt > 1 {
		color := r.light.Average(r.cursor).Average()
		// TODO: write a real function that returns the ideal sample count weighting midtones and shadows more heavily than blown-out highlights
		midtone := 200 - math.Min(math.Abs(25-color), 200)
		multiplier := midtone / 200
		noise := r.light.Noise(r.cursor)
		v, c := r.light.Variance()
		noiseRatio := (noise + 1) / (v + 1)
		targetCount := noiseRatio * c * multiplier
		correction := (targetCount - r.light.Count(r.cursor))
		adapted := math.Max(0, math.Min(r.adapt, correction))
		count += uint(adapted)
	}
	for i := uint(0); i < count; i++ {
		pixels <- r.cursor
	}
	r.cursor = (r.cursor + 1) % uint(r.width*r.height)
	if r.cursor == 0 {
		r.light.UpdateVariance()
	}
	return count
}

func (r *Render) WritePngs(out, heat, noise string, expose float64) error {
	if err := writePNG(out, r.Light().Image(expose)); err != nil {
		return err
	}
	if len(heat) > 0 {
		if err := writePNG(heat, r.Light().Heatmap()); err != nil {
			return err
		}
	}
	if len(noise) > 0 {
		if err := writePNG(noise, r.Light().Noisemap()); err != nil {
			return err
		}
	}
	return nil
}

func writePNG(file string, i image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, i)
	return err
}
