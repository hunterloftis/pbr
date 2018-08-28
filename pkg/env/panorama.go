package env

import (
	"errors"
	"math"
	"os"

	"github.com/Opioid/rgbe"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

const maxEnergy = 1000000

type Pano struct {
	Expose float64
	width  int
	height int
	data   []float32
}

// http://gl.ict.usc.edu/Data/HighResProbes/
func (p *Pano) At(dir geom.Dir) rgb.Energy {
	u := 1 + math.Atan2(dir.X, -dir.Z)/math.Pi // [0,2]
	v := math.Acos(dir.Y) / math.Pi            // [0,1]
	x := int(u / 2 * float64(p.width))
	y := int(v * float64(p.height))
	i := ((y*p.width + x) * 3) % len(p.data)
	energy := rgb.Energy{
		X: float64(p.data[i]),
		Y: float64(p.data[i+1]),
		Z: float64(p.data[i+2]),
	}
	return energy.Scaled(p.Expose).Limit(maxEnergy)
}

func ReadFile(filename string, expose float64) (*Pano, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	width, height, data, err := rgbe.Decode(f)
	if err != nil {
		return nil, err
	}
	if width/height != 2 {
		return nil, errors.New("Unsupported HDRI dimensions (need 2:1 aspect ratio)")
	}
	p := Pano{
		Expose: expose,
		width:  width,
		height: height,
		data:   data,
	}
	return &p, nil
}
