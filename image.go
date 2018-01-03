package pbr

import (
	"image"
	"math"
)

// Pixel elements are stored in specific offsets.
// These constants allow easy access, eg `someFloat64Array[i + Blue]`
const (
	Red uint = iota
	Green
	Blue
	Count
	Noise
	Stride
)

type Image struct {
	Width, Height uint
	pixels        []float64 // stored in a flat array, chunked by Stride
	variance      float64
}

func NewImage(width, height uint) Image {
	return Image{
		Width:  width,
		Height: height,
		pixels: make([]float64, width*height*Stride),
	}
}

// TODO: this is confusing. Currently is the size of the array, but could also be width * height
func (im Image) Size() uint {
	return uint(len(im.pixels))
}

func (im Image) Color(i uint) (r, g, b float64) {
	return im.pixels[i+Red], im.pixels[i+Green], im.pixels[i+Blue]
}

func (im Image) Count(i uint) float64 {
	return im.pixels[i+Count]
}

// TODO: these should all factor in stride and take "i" as a pixel index, vs a slice index
func (im Image) Noise(i uint) float64 {
	return im.pixels[i+Noise]
}

func (im Image) Average(pixel uint) Energy {
	i := pixel * Stride
	c := float64(im.pixels[i+Count])
	red := im.pixels[i+Red] / c
	green := im.pixels[i+Green] / c
	blue := im.pixels[i+Blue] / c
	return Energy{red, green, blue}
}

// Rgb averages each sample into an rgb value.
func (im Image) Rgb(expose float64) image.Image {
	rgba := image.NewRGBA(image.Rect(0, 0, int(im.Width), int(im.Height)))
	length := im.Size()
	for i := uint(0); i < length; i += Stride {
		i2 := i / Stride * 4
		count := im.Count(i)
		r, g, b := im.Color(i)
		rgba.Pix[i2] = color(r / count * expose)
		rgba.Pix[i2+1] = color(g / count * expose)
		rgba.Pix[i2+2] = color(b / count * expose)
		rgba.Pix[i2+3] = 255
	}
	return rgba
}

// Heat returns a heatmap of the sample count for each pixel.
func (im Image) Heat(offset uint) image.Image {
	rgba := image.NewRGBA(image.Rect(0, 0, int(im.Width), int(im.Height)))
	max := 0.0
	length := im.Size()
	for i := uint(0); i < length; i += Stride {
		max = math.Max(max, im.pixels[i+offset])
	}
	for i := uint(0); i < length; i += Stride {
		i2 := i / Stride * 4
		c := color(im.pixels[i+offset] / max * 255)
		rgba.Pix[i2] = c
		rgba.Pix[i2+1] = c
		rgba.Pix[i2+2] = c
		rgba.Pix[i2+3] = 255
	}
	return rgba
}

func (im Image) Integrate(index uint, sample Energy, noise bool) {
	p := index * Stride
	rgb := [3]float64{sample.X, sample.Y, sample.Z}
	im.pixels[p+Red] += rgb[0]
	im.pixels[p+Green] += rgb[1]
	im.pixels[p+Blue] += rgb[2]
	im.pixels[p+Count]++
	if noise {
		p := index * Stride
		mean := im.Average(index)
		variance := sample.Variance(mean)
		count := im.pixels[p+Count]
		oldNoise := im.pixels[p+Noise] * (count - 1) / count
		newNoise := variance / count
		im.pixels[p+Noise] = oldNoise + newNoise
	}
}

func (im Image) UpdateVariance() float64 {
	size := im.Size()
	count := float64(im.Width * im.Height)
	im.variance = 0
	for i := uint(0); i < size; i += Stride {
		im.variance += im.pixels[i+Noise] / count
	}
	return im.variance
}

func (im Image) Variance() float64 {
	return im.variance
}

func color(n float64) uint8 {
	return uint8(gamma(math.Min(n, 255), 2.2))
}

func gamma(n, g float64) float64 {
	return math.Pow(n/255, (1/g)) * 255
}
