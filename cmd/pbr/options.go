package main

import (
	"math"
	"path/filepath"

	arg "github.com/alexflint/go-arg"
	"github.com/hunterloftis/pbr"
)

// Options configures rendering behavior.
type Options struct {
	Scene string `arg:"positional,required" help:"input scene .obj"`
	Info  bool   `help:"output scene information and exit"`

	Out     string `help:"output render .png"`
	Heat    string `help:"output heatmap as .png"`
	Noise   string `help:"output noisemap as .png"`
	Profile bool   `help:"record performance into profile.pprof"`

	Width  int `help:"rendering width in pixels"`
	Height int `help:"rendering height in pixels"`

	Sky    *pbr.Energy `help:"ambient sky color"`
	Ground *pbr.Energy `help:"ground color"`
	Env    string      `help:"environment as a panoramic hdr radiosity map (.hdr file)"`
	Rad    float64     `help:"exposure of the hdr (radiosity) environment map"`

	Adapt    float64 `help:"adaptive sampling multiplier"`
	Bounce   int     `help:"number of light bounces"`
	Direct   int     `help:"number of direct rays to cast"`   // TODO: implement
	Indirect int     `help:"number of indirect rays to cast"` // TODO: implement
	Complete float64 `help:"number of samples-per-pixel at which to exit"`

	From   *pbr.Vector3 `help:"camera position"`
	To     *pbr.Vector3 `help:"camera target"`
	Focus  *pbr.Vector3 `help:"camera focus (if other than 'to')"`
	Lens   float64      `help:"camera focal length in mm"`
	FStop  float64      `help:"camera f-stop"`
	Expose float64      `help:"exposure multiplier"`
}

func options() *Options {
	c := &Options{
		Width:    800,
		Height:   600,
		Profile:  false,
		Sky:      &pbr.Energy{210, 230, 255},
		Ground:   &pbr.Energy{0, 0, 0},
		Rad:      100,
		Adapt:    10,
		Bounce:   10,
		Direct:   1,
		Indirect: 1,
		Complete: math.Inf(1),
		Lens:     50,
		FStop:    4,
		Expose:   1,
	}
	arg.MustParse(c)
	if c.Out == "" && !c.Info {
		name := filepath.Base(c.Scene)
		ext := filepath.Ext(name)
		c.Out = name[0:len(name)-len(ext)] + ".png"
	}
	return c
}

func cameraOptions(o *Options, bounds *pbr.Box, center pbr.Vector3) (from, to, focus pbr.Vector3) {
	if o.From == nil {
		twoThirds := pbr.Vector3{bounds.Max.X * 9, bounds.Max.Y, bounds.Max.Z * 6}
		from = twoThirds
	}
	if o.To == nil {
		to = center
	}
	if o.Focus == nil {
		focus = to
	}
	return from, to, focus
}
