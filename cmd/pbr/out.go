package main

import (
	"fmt"
	"os"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/geom"
)

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
}

func printInfo(b *geom.Bounds, surfaces int, c *camera.SLR) {
	fmt.Println("Min:", b.Min)
	fmt.Println("Max:", b.Max)
	fmt.Println("Center:", b.Center)
	fmt.Println("Camera:", c)
}
