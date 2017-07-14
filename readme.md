# pbr: a Physically-Based Renderer in Go

Package pbr implements Physically-Based Rendering with a Monte Carlo path tracer.
[[ Documentation ]](https://godoc.org/github.com/hunterloftis/pbr/pbr)
[[ Github ]](https://github.com/hunterloftis/pbr)

[![Render](https://user-images.githubusercontent.com/364501/27998485-ff50ece2-64dd-11e7-861b-a8fb336d6e50.png)](https://user-images.githubusercontent.com/364501/27998485-ff50ece2-64dd-11e7-861b-a8fb336d6e50.png)

- Unbiased [Monte-Carlo integration](https://en.wikipedia.org/wiki/Monte_Carlo_integration)
- Adaptive [sampling](https://renderman.pixar.com/resources/RenderMan_20/risSampling.html)
- [Russian roulette](https://computergraphics.stackexchange.com/questions/2316/is-russian-roulette-really-the-answer)
- Parametric shapes (spheres, cubes)
- Transformation matrices (translate, rotate, scale)
- Physically-based [materials](https://www.marmoset.co/posts/basic-theory-of-physically-based-rendering/)
  - Fresnel reflection, transmission, absorption, diffusion
  - Color, refractive indices, gloss, transparency, separate fresnel channels, metals
- Arbitrary light sources ('everything is a light')
- [Environment maps](http://gl.ict.usc.edu/Data/HighResProbes/)
- Physically-based cameras
  - Sensor, aperture, focus, depth-of-field
- [Supersampled anti-aliasing](https://en.wikipedia.org/wiki/Supersampling)
- Fully concurrent with a sequential API

## Hello, world

```
$ go get -u github.com/hunterloftis/pbr/pbr
$ cd $GOPATH/src/github.com/hunterloftis/pbr
$ ./hello
```

![Hello, world render](https://user-images.githubusercontent.com/364501/28223346-111a1944-6899-11e7-946b-8ea5c90c3888.png)

```go
func main() {
	scene := pbr.EmptyScene()
	camera := pbr.NewCamera(960, 540)
	sampler := pbr.NewSampler(camera, scene)
	renderer := pbr.NewRenderer(sampler)

	scene.SetSky(pbr.Vector3{256, 256, 256}, pbr.Vector3{})
	scene.Add(pbr.UnitSphere(pbr.Plastic(1, 0, 0, 1)))

	for sampler.PerPixel() < 200 {
		sampler.Sample()
		fmt.Printf("\r%.1f samples / pixel", sampler.PerPixel())
	}
	pbr.WritePNG("hello.png", renderer.Rgb())
}
```

## Other examples

```
$ ./cubes
$ ./render
```

## Testing

```
$ go test ./pbr
```