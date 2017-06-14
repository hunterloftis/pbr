package main

import (
	"github.com/hunterloftis/trace/trace"
)

func main() {
	renderer := trace.Renderer{Width: 960, Height: 540}
	renderer.Write("test.png")

}
