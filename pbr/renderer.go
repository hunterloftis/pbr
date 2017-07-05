package pbr

import (
	"image"
	"image/png"
	"math"
	"os"
)

// Renderer renders the results of a trace to a file
type Renderer struct {
	Width    int
	Height   int
	Exposure float64
	pixels   []float64
}

// CamRenderer sizes a Renderer to match a Camera
func CamRenderer(cam *Camera, exp float64) *Renderer {
	return &Renderer{
		Width:    cam.Width,
		Height:   cam.Height,
		Exposure: exp,
		pixels:   make([]float64, 0),
	}
}

// Merge merges pixel arrays
func (r *Renderer) Merge(pixels []float64) {
	if len(r.pixels) < len(pixels) {
		r.pixels = make([]float64, len(pixels))
	}
	for i, val := range pixels {
		r.pixels[i] += val
	}
}

// Rgb averages each sample into an rgb value
func (r *Renderer) Rgb() []float64 {
	rgb := make([]float64, len(r.pixels)/PROPS*3)
	for i := 0; i < len(r.pixels); i += PROPS {
		i2 := i / PROPS * 3
		count := r.pixels[i+3]
		rgb[i2] = r.pixels[i] / count
		rgb[i2+1] = r.pixels[i+1] / count
		rgb[i2+2] = r.pixels[i+2] / count
	}
	return rgb
}

// Heat returns a heatmap of the sample count for each pixel
func (r *Renderer) Heat() []float64 {
	heat := make([]float64, len(r.pixels)/PROPS*3)
	max := 0.0
	for i := 3; i < len(r.pixels); i += PROPS {
		max = math.Max(max, r.pixels[i])
	}
	for i := 0; i < len(r.pixels); i += PROPS {
		i2 := i / PROPS * 3
		heat[i2] = r.pixels[i+3] / max * 255
		heat[i2+1] = r.pixels[i+3] / max * 255
		heat[i2+2] = r.pixels[i+3] / max * 255
	}
	return heat
}

// WriteRGB writes RGB data to a file
func (r *Renderer) WriteRGB(file string) error {
	return r.Write(r.Rgb(), file)
}

// WriteHeat writes heat (count) data to a file
func (r *Renderer) WriteHeat(file string) error {
	return r.Write(r.Heat(), file)
}

// Write writes RGB data to a png
func (r *Renderer) Write(rgb []float64, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	m := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	for i := 0; i < len(rgb); i += 3 {
		i2 := i / 3 * 4
		m.Pix[i2] = r.color(rgb[i])
		m.Pix[i2+1] = r.color(rgb[i+1])
		m.Pix[i2+2] = r.color(rgb[i+2])
		m.Pix[i2+3] = 255
	}
	return png.Encode(f, m)
}

func (r *Renderer) color(n float64) uint8 {
	return uint8(gamma(math.Min(n*r.Exposure, 255), 2.2))
}

func gamma(n, g float64) float64 {
	return math.Pow(n/255, (1/g)) * 255
}
