package main

import (
	"github.com/hunterloftis/trace/trace"
)

func main() {
	sampler := trace.NewSampler(960, 540)
	renderer := trace.Renderer{Width: 960, Height: 540}
	for i := 0; i < 10; i++ {
		sampler.Sample()
	}
	renderer.Write(sampler.Values(), "test.png")
}
