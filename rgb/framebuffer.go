package rgb

import (
	"image"
	"math"
)

// TODO: enforce immutability by creating an interface that exposes only readable properties

// Pixel elements are stored in specific offsets.
// These constants allow easy access, eg `someFloat64Array[i + blue]`
const (
	red uint = iota
	green
	blue
	count
	noise
	stride
)

// TODO: rename to Framebuffer
type Framebuffer struct {
	width, height uint
	meanVariance  float64
	meanCount     float64
	pixels        []float64 // stored in a flat array, chunked by stride
}

func NewBuffer(width, height uint) *Framebuffer {
	return &Framebuffer{
		width:  width,
		height: height,
		pixels: make([]float64, width*height*stride),
	}
}

func (f *Framebuffer) Count(pixel uint) float64 {
	return f.pixels[(pixel*stride)+count]
}

func (f *Framebuffer) Noise(pixel uint) float64 {
	return f.pixels[(pixel*stride)+noise]
}

func (f *Framebuffer) Average(pixel uint) Energy {
	i := pixel * stride
	c := float64(f.pixels[i+count]) + 1e-10 // TODO: is this the best way to avoid a divide by zero?
	r := f.pixels[i+red] / c
	g := f.pixels[i+green] / c
	b := f.pixels[i+blue] / c
	return Energy{r, g, b}
}

// Rgb averages each sample into an rgb value.
// https://stackoverflow.com/questions/21984405/relation-between-sigma-and-radius-on-the-gaussian-blur
func (f *Framebuffer) Image(expose float64) image.Image {
	rgba := image.NewRGBA(image.Rect(0, 0, int(f.width), int(f.height)))
	size := f.size()
	for i := uint(0); i < size; i++ {
		i2 := i * 4
		color := f.Average(i)
		rgba.Pix[i2] = tonemap(color.X * expose)
		rgba.Pix[i2+1] = tonemap(color.Y * expose)
		rgba.Pix[i2+2] = tonemap(color.Z * expose)
		rgba.Pix[i2+3] = 255
	}
	return rgba
}

func (f *Framebuffer) Heatmap() image.Image {
	return f.heat(count)
}

func (f *Framebuffer) Noisemap() image.Image {
	return f.heat(noise)
}

// Heat returns a heatmap of the sample count for each pixel.
func (f *Framebuffer) heat(offset uint) image.Image {
	rgba := image.NewRGBA(image.Rect(0, 0, int(f.width), int(f.height)))
	max := 0.0
	size := f.size()
	for i := uint(0); i < size; i++ {
		max = math.Max(max, f.val(i, offset))
	}
	for i := uint(0); i < size; i++ {
		i2 := i * 4
		c := tonemap(f.val(i, offset) / max * 255)
		rgba.Pix[i2] = c
		rgba.Pix[i2+1] = c
		rgba.Pix[i2+2] = c
		rgba.Pix[i2+3] = 255
	}
	return rgba
}

func (f *Framebuffer) Integrate(index uint, sample Energy) uint {
	p := index * stride
	rgb := [3]float64{sample.X, sample.Y, sample.Z}
	f.pixels[p+red] += rgb[0]
	f.pixels[p+green] += rgb[1]
	f.pixels[p+blue] += rgb[2]
	f.pixels[p+count]++

	// noise
	mean := f.Average(index)
	variance := (sample.Variance(mean) + 1) / (mean.Average() + 1)
	c := f.pixels[p+count]
	oldNoise := f.pixels[p+noise] * (c - 1) / c
	newNoise := variance / c
	f.pixels[p+noise] = oldNoise + newNoise
	return 1
}

func (f *Framebuffer) UpdateVariance() {
	f.meanVariance = 0
	f.meanCount = 0
	size := f.size()
	for i := uint(0); i < size; i++ {
		f.meanVariance += f.val(i, noise) / float64(size)
		f.meanCount += f.val(i, count) / float64(size)
	}
}

func (f *Framebuffer) Variance() (v, c float64) {
	return f.meanVariance, f.meanCount
}

func (f *Framebuffer) val(pixel, offset uint) float64 {
	return f.pixels[(pixel*stride)+offset]
}

func (f *Framebuffer) size() uint {
	return f.width * f.height
}

func tonemap(n float64) uint8 {
	return uint8(gamma(math.Min(n, 255), 2.2))
}

func gamma(n, g float64) float64 {
	return math.Pow(n/255, (1/g)) * 255
}

func (f *Framebuffer) bilateral(pixel uint, weightS, weightR float64, radius int) Energy {
	e0 := f.Average(pixel).Limit(255)
	// m0 := f.material(pixel)
	x0, y0 := int(pixel%f.width), int(pixel/f.width)
	eSum := Energy{}
	norm := 0.0
	for y := y0 - radius; y <= y0+radius; y++ {
		for x := x0 - radius; x <= x0+radius; x++ {
			if y < 0 || y >= int(f.height) || x < 0 || x >= int(f.width) {
				continue
			}
			i := uint(y)*f.width + uint(x)
			e := f.Average(i).Limit(255)
			dx := x - x0
			dy := y - y0
			s := math.Sqrt(float64(dx*dx + dy*dy))
			r := e.Minus(e0).Size() / 442 // size of Energy{255, 255, 255}
			weight := gaussian(s, weightS) * gaussian(r, weightR)
			// if f.material(i) != m0 {
			// 	weight *= 0.3
			// }
			eSum = eSum.Plus(e.Amplified(weight))
			norm += weight
		}
	}
	return eSum.Amplified(1 / norm)
}

// https://dsp.stackexchange.com/questions/10057/gaussian-blur-standard-deviation-radius-and-kernel-size
func gaussian(n, sigma float64) float64 {
	return math.Exp(-(n*n)/(2*sigma*sigma)) / (sigma * math.Sqrt(2*math.Pi))
}
