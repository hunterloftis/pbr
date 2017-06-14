package main

import (
	"github.com/hunterloftis/trace/trace"
)

func main() {
	sampler := trace.Sampler{Width: 960, Height: 540}
	renderer := trace.Renderer{Width: 960, Height: 540}
	renderer.Write(sampler.Samples(), "test.png")
}
