package main

import (
	"math"
	"path/filepath"

	arg "github.com/alexflint/go-arg"
	"github.com/hunterloftis/pbr"
)

// Options configures rendering behavior.
// TODO: add "watermark"
// TODO: --filter (per pixel samples after which to apply smoothing filters; 0 = off)
type Options struct {
	Scene string `arg:"positional,required" help:"input scene .obj"`
	Info  bool   `help:"output scene information and exit"`

	Out     string `help:"output render .png"`
	Heat    string `help:"output heatmap as .png"`
	Noise   string `help:"output noisemap as .png"`
	Profile bool   `help:"record performance into profile.pprof"`

	Width  int `help:"rendering width in pixels"`
	Height int `help:"rendering height in pixels"`

	Sky      *pbr.Energy `help:"ambient sky color"`
	Ground   *pbr.Energy `help:"ground color"`
	Env      string      `help:"environment as a panoramic hdr radiosity map (.hdr file)"`
	Rad      float64     `help:"exposure of the hdr (radiosity) environment map"`
	Floor    bool        `help:"create a floor underneath the scene"`
	Adapt    float64     `help:"adaptive sampling multiplier"`
	Bounce   int         `help:"number of light bounces"`
	Direct   int         `help:"number of direct rays to cast"`   // TODO: implement
	Indirect int         `help:"number of indirect rays to cast"` // TODO: implement
	Complete float64     `help:"number of samples-per-pixel at which to exit"`
	Thin     bool        `help:"treat transparent surfaces as having zero thickness"`

	From      *pbr.Vector3 `help:"camera position"`
	To        *pbr.Vector3 `help:"camera target"`
	Focus     *pbr.Vector3 `help:"camera focus (if other than 'to')"`
	Dist      float64      `help:"camera distance from target"`
	Polar     float64      `help:"camera polar angle on target"`
	Longitude float64      `help:"camera longitudinal angle on target"`
	Lens      float64      `help:"camera focal length in mm"`
	FStop     float64      `help:"camera f-stop"`
	Expose    float64      `help:"exposure multiplier"`
}

func options() *Options {
	c := &Options{
		Width:     800,
		Height:    600,
		Profile:   false,
		Sky:       &pbr.Energy{210, 230, 255},
		Ground:    &pbr.Energy{0, 0, 0},
		Rad:       100,
		Adapt:     10,
		Bounce:    8,
		Direct:    1,
		Indirect:  1,
		Complete:  math.Inf(1),
		Lens:      50,
		Polar:     0,
		Longitude: 0,
		FStop:     4,
		Expose:    1,
	}
	arg.MustParse(c)
	if c.Out == "" && !c.Info {
		name := filepath.Base(c.Scene)
		ext := filepath.Ext(name)
		c.Out = name[0:len(name)-len(ext)] + ".png"
	}
	return c
}

// TODO: accept input that configures this smart behavior instead of specifying exact camera points
// eg, -distance 5 -theta 0 -phi 3.14
func completeOptions(o *Options, bounds *pbr.Box, center pbr.Vector3, surfaces []pbr.Surface) {
	if o.From == nil {
		dir := pbr.AngleDirection(o.Polar, o.Longitude)
		ray := pbr.NewRay(center, dir)
		theta := pbr.FieldOfView(o.Lens, 35) / 2
		if o.Dist == 0 {
			max := o.Dist
			for _, s := range surfaces {
				pt := s.Center()
				dist0 := pt.Minus(center).Dot(pbr.Vector3(dir))
				offset := pt.Minus(ray.Moved(dist0)).Len()
				dist1 := offset / math.Tan(theta)
				dist := dist0 + dist1
				if dist > max {
					max = dist
				}
			}
			o.Dist = max * 1.1
		}
		from := ray.Moved(o.Dist)
		o.From = &from
	}
	if o.To == nil {
		o.To = &center
	}
	if o.Focus == nil {
		o.Focus = o.To
	}
}
