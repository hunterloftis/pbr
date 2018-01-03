package main

import (
	"math"

	arg "github.com/alexflint/go-arg"
	"github.com/hunterloftis/pbr"
)

// Options configures rendering behavior.
type Options struct {
	Scene   string `arg:"positional,required" help:"input scene .obj"`
	Render  string `arg:"positional,required" help:"output render .png"`
	Heat    string `help:"output heatmap as .png"`
	Noise   string `help:"output noisemap as .png"`
	Profile bool   `help:"record performance into profile.pprof"`
	Width   int    `help:"rendering width in pixels"`
	Height  int    `help:"rendering height in pixels"`

	Env    string      `help:"environment as a panoramic hdr radiosity map"`
	Sky    *pbr.Energy `help:"ambient sky color"`
	Ground *pbr.Energy `help:"ground color"`

	Adapt    float64 `help:"adaptive sampling multiplier"`
	Bounce   int     `help:"number of light bounces"`
	Direct   int     `help:"number of direct rays to cast"`
	Indirect int     `help:"number of indirect rays to cast"`
	Complete float64 `help:"number of samples-per-pixel at which to exit"`

	From   *pbr.Vector3 `help:"camera position"`
	To     *pbr.Vector3 `help:"camera target"`
	Focus  *pbr.Vector3 `help:"camera focus (if other than 'to')"`
	Lens   float64      `help:"camera focal length in mm"`
	FStop  float64      `help:"camera f-stop"`
	Expose float64      `help:"exposure multiplier"`
}

func options() Options {
	c := Options{
		Width:    1280,
		Height:   720,
		Profile:  false,
		Sky:      &pbr.Energy{40, 50, 60},
		Ground:   &pbr.Energy{0, 0, 0},
		Adapt:    10,
		Bounce:   10,
		Direct:   1,
		Indirect: 1,
		Complete: math.Inf(1),
		From:     &pbr.Vector3{0, 0, 2},
		To:       &pbr.Vector3{0, 0, 0},
		Focus:    nil,
		Lens:     50,
		FStop:    4,
		Expose:   1,
	}
	arg.MustParse(&c)
	if c.Focus == nil {
		c.Focus = c.To
	}
	return c
}
