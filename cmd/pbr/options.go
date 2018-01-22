package main

import (
	"math"
	"path/filepath"

	arg "github.com/alexflint/go-arg"
	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

// Options configures rendering behavior.
// TODO: add "watermark"
// TODO: --filter (per pixel samples after which to apply smoothing filters; 0 = off)
// TODO: change polar/latitude to lat/lon
type Options struct {
	Scene   string `arg:"positional,required" help:"input scene .obj"`
	Verbose bool   `help:"verbose output with scene information"`
	Info    bool   `arg:"-i" help:"output scene information and exit"`

	Out     string `arg:"-o" help:"output render .png"`
	Heat    string `help:"output heatmap as .png"`
	Noise   string `help:"output noisemap as .png"`
	Profile bool   `help:"record performance into profile.pprof"`

	Width  int `help:"rendering width in pixels"`
	Height int `help:"rendering height in pixels"`

	Target *geom.Vector3 `help:"camera target point"`
	Focus  *geom.Vector3 `help:"camera focus point (if other than 'target')"`
	Dist   float64       `help:"camera distance from target"`
	Lat    float64       `help:"camera polar angle on target"`
	Lon    float64       `help:"camera longitudinal angle on target"`

	Lens   float64 `help:"camera focal length in mm"`
	FStop  float64 `help:"camera f-stop"`
	Expose float64 `help:"exposure multiplier"`

	Ambient  *rgb.Energy `help:"the ambient light color"`
	Env      string      `arg:"-e" help:"environment as a panoramic hdr radiosity map (.hdr file)"`
	Rad      float64     `help:"exposure of the hdr (radiosity) environment map"`
	Floor    bool        `help:"create a floor underneath the scene"`
	Adapt    float64     `help:"adaptive sampling multiplier"`
	Bounce   int         `arg:"-d" help:"number of light bounces (depth)"`
	Direct   int         `arg:"-d" help:"maximum number of direct rays to cast"`
	Branch   int         `arg:"-b" help:"maximum number of branches on first hit"`
	Complete float64     `arg:"-c" help:"number of samples-per-pixel at which to exit"`
	Thin     bool        `help:"treat transparent surfaces as having zero thickness"`
}

func options() *Options {
	c := &Options{
		Width:    800,
		Height:   450,
		Ambient:  &rgb.Energy{500, 500, 500},
		Rad:      100,
		Adapt:    8,
		Bounce:   8,
		Direct:   1,
		Branch:   32,
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

func (o *Options) Version() string {
	return "1.0.0"
}

func (o *Options) Description() string {
	return "pbr renders physically-based 3D scenes with a Monte Carlo path tracer."
}
