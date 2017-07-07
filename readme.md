# pbr: a Physically-Based Renderer in Go

Package pbr implements Physically-Based Rendering with a Monte Carlo path tracer.
[[ API Docs ]](https://godoc.org/github.com/hunterloftis/pbr/pbr)
[[ Source ]](https://github.com/hunterloftis/pbr)

![Render](https://user-images.githubusercontent.com/364501/27976627-5e28c4ba-6335-11e7-9fde-00c7ddeb5987.png)

- Unbiased [Monte-Carlo integration](https://en.wikipedia.org/wiki/Monte_Carlo_integration)
- Adaptive [sampling](https://renderman.pixar.com/resources/RenderMan_20/risSampling.html)
- [Russian roulette](https://computergraphics.stackexchange.com/questions/2316/is-russian-roulette-really-the-answer)
- Parametric shapes (spheres, cubes)
- Physically-based [materials](https://www.marmoset.co/posts/basic-theory-of-physically-based-rendering/)
  - Fresnel reflection, transmission, absorption, diffusion
  - Color, refractive indices, gloss, transparency, separate fresnel channels, metals
- Arbitrary light sources ('everything is a light')
- [Environment maps](http://gl.ict.usc.edu/Data/HighResProbes/)
- Physically-based cameras
  - Sensor, aperture, focus, depth-of-field
- [Supersampled anti-aliasing](https://en.wikipedia.org/wiki/Supersampling)

## Quick start

```
$ go get -u github.com/hunterloftis/pbr/pbr
$ cd $GOPATH/src/github.com/hunterloftis/pbr
$ ./run
```

By default, the renderer runs until it receives a signal (like Ctrl + C)

## Scene bins

Scenes (like the example `cmd/cubes.go` scene) are built into binaries:

```
$ go build -o bin/cubes ./cmd/cubes.go
$ bin/cubes -help
Usage of bin/cubes:
  -adapt int
    	Adaptive sampling; 0=off, 3=medium, 5=high (default 4)
  -bounces int
    	Maximum light bounces (default 10)
  -heat string
    	Heatmap png filename
  -out string
    	Output png filename (default "render.png")
  -profile
    	Record performance into profile.pprof
  -samples float
    	Max samples per pixel (default +Inf)
  -workers int
    	Concurrency level (default 4)
```

## Testing

```
$ go test ./pbr
```