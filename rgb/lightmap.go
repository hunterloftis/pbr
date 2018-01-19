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
type Lightmap struct {
	width, height uint
	meanVariance  float64
	meanCount     float64
	pixels        []float64 // stored in a flat array, chunked by stride
}

func Map(width, height uint) *Lightmap {
	return &Lightmap{
		width:  width,
		height: height,
		pixels: make([]float64, width*height*stride),
	}
}

func (lm *Lightmap) Count(pixel uint) float64 {
	return lm.pixels[(pixel*stride)+count]
}

func (lm *Lightmap) Noise(pixel uint) float64 {
	return lm.pixels[(pixel*stride)+noise]
}

func (lm *Lightmap) Average(pixel uint) Energy {
	i := pixel * stride
	c := float64(lm.pixels[i+count]) + 1e-10 // TODO: is this the best way to avoid a divide by zero?
	r := lm.pixels[i+red] / c
	g := lm.pixels[i+green] / c
	b := lm.pixels[i+blue] / c
	return Energy{r, g, b}
}

// Rgb averages each sample into an rgb value.
func (lm *Lightmap) Image(expose float64) image.Image {
	rgba := image.NewRGBA(image.Rect(0, 0, int(lm.width), int(lm.height)))
	size := lm.size()
	for i := uint(0); i < size; i++ {
		i2 := i * 4
		c := lm.Count(i)
		r, g, b := lm.color(i)
		rgba.Pix[i2] = tonemap(r / c * expose)
		rgba.Pix[i2+1] = tonemap(g / c * expose)
		rgba.Pix[i2+2] = tonemap(b / c * expose)
		rgba.Pix[i2+3] = 255
	}
	return rgba
}

func (lm *Lightmap) Heatmap() image.Image {
	return lm.heat(count)
}

func (lm *Lightmap) Noisemap() image.Image {
	return lm.heat(noise)
}

// Heat returns a heatmap of the sample count for each pixel.
func (lm *Lightmap) heat(offset uint) image.Image {
	rgba := image.NewRGBA(image.Rect(0, 0, int(lm.width), int(lm.height)))
	max := 0.0
	size := lm.size()
	for i := uint(0); i < size; i++ {
		max = math.Max(max, lm.val(i, offset))
	}
	for i := uint(0); i < size; i++ {
		i2 := i * 4
		c := tonemap(lm.val(i, offset) / max * 255)
		rgba.Pix[i2] = c
		rgba.Pix[i2+1] = c
		rgba.Pix[i2+2] = c
		rgba.Pix[i2+3] = 255
	}
	return rgba
}

func (lm *Lightmap) Integrate(index uint, sample Energy) {
	p := index * stride
	rgb := [3]float64{sample.X, sample.Y, sample.Z}
	lm.pixels[p+red] += rgb[0]
	lm.pixels[p+green] += rgb[1]
	lm.pixels[p+blue] += rgb[2]
	lm.pixels[p+count]++

	// noise
	// TODO: profile
	mean := lm.Average(index)
	variance := sample.Variance(mean)
	c := lm.pixels[p+count]
	oldNoise := lm.pixels[p+noise] * (c - 1) / c
	newNoise := variance / c
	lm.pixels[p+noise] = oldNoise + newNoise
}

func (lm *Lightmap) UpdateVariance() {
	lm.meanVariance = 0
	lm.meanCount = 0
	size := lm.size()
	for i := uint(0); i < size; i++ {
		lm.meanVariance += lm.val(i, noise) / float64(size)
		lm.meanCount += lm.val(i, count) / float64(size)
	}
}

func (lm *Lightmap) Variance() (v, c float64) {
	return lm.meanVariance, lm.meanCount
}

func (lm *Lightmap) val(pixel, offset uint) float64 {
	return lm.pixels[(pixel*stride)+offset]
}

func (lm *Lightmap) size() uint {
	return lm.width * lm.height
}

func (lm *Lightmap) color(pixel uint) (r, g, b float64) {
	i := pixel * stride
	return lm.pixels[i+red], lm.pixels[i+green], lm.pixels[i+blue]
}

func tonemap(n float64) uint8 {
	return uint8(gamma(math.Min(n, 255), 2.2))
}

func gamma(n, g float64) float64 {
	return math.Pow(n/255, (1/g)) * 255
}
