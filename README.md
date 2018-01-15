# pbr: a Physically-Based 3D Renderer in Go

Package pbr implements Physically-Based Rendering with a Monte Carlo path tracer.
[[ Documentation ]](https://godoc.org/github.com/hunterloftis/pbr)
[[ Github ]](https://github.com/hunterloftis/pbr)

![Render](https://user-images.githubusercontent.com/364501/34923521-c39b132c-f96a-11e7-9a27-f79f67268079.png)

```
$ pbr lambo.obj -floor -polar 3.6 -longitude 0.1 -env 293.hdr -rad 500 -lens 60 -fstop 1.4 -to=-0.1,0.5,0.1 -dist 7.5 -focus=-2.2658,0.5542,0
```

- [Unbiased Monte-Carlo integration](https://en.wikipedia.org/wiki/Monte_Carlo_integration)
- [Adaptive sampling](https://renderman.pixar.com/resources/RenderMan_20/risSampling.html)
- [Russian roulette](https://computergraphics.stackexchange.com/questions/2316/is-russian-roulette-really-the-answer)
- Parametric shapes (spheres, cubes, triangles)
- Transformation matrices (translate, rotate, scale)
- Wavefront .obj files (meshes) and .mtl files (materials)
	- [With extended physically-based material properties](http://exocortex.com/blog/extending_wavefront_mtl_to_support_pbr)
- [Physically-based materials](https://www.marmoset.co/posts/basic-theory-of-physically-based-rendering/)
  - Fresnel reflection, transmission, absorption, diffusion
  - Color, refractive indices, gloss, transparency, separate fresnel channels, metals
- Arbitrary light sources ('everything is a light')
- [Environment maps](http://gl.ict.usc.edu/Data/HighResProbes/)
	- [Image-based lighting](https://agraphicsguy.wordpress.com/2016/09/07/image-based-lighting-in-offline-and-real-time-rendering/)
- Physically-based cameras
  - Sensor, aperture, focal length, focus, depth-of-field
- [Supersampled anti-aliasing](https://en.wikipedia.org/wiki/Supersampling)
- 100% Go with no system dependencies
	- Sequential API, concurrent execution
	- CLI

## Try it

Download:

```
$ go get github.com/hunterloftis/pbr
$ cd $GOPATH/src/github.com/hunterloftis/pbr
$ dep ensure
```

Run:

```
$ go install github.com/hunterloftis/pbr/cmd/pbr
$ pbr $GOPATH/src/github.com/hunterloftis/pbr/fixtures/models/falcon.obj -longitude 0.5 -polar 0.5 -complete 5
$ open falcon.png
```

![falcon render](https://user-images.githubusercontent.com/364501/34923876-aaa1f07c-f96d-11e7-8dc2-199373bdb51e.png)

## API

See [GoDoc](https://godoc.org/github.com/hunterloftis/pbr)

## CLI

```
Usage: pbr SCENE [options]

Positional arguments:
  SCENE                  input scene .obj

Options:
  --info                 output scene information and exit
  --out OUT              output render .png
  --heat HEAT            output heatmap as .png
  --noise NOISE          output noisemap as .png
  --profile              record performance into profile.pprof
  --width WIDTH          rendering width in pixels [default: 800]
  --height HEIGHT        rendering height in pixels [default: 600]
  --sky SKY              ambient sky color [default: &{210 230 255}]
  --ground GROUND        ground color [default: &{0 0 0}]
  --env ENV              environment as a panoramic hdr radiosity map (.hdr file)
  --rad RAD              exposure of the hdr (radiosity) environment map [default: 100]
  --floor                create a floor underneath the scene
  --adapt ADAPT          adaptive sampling multiplier [default: 10]
  --bounce BOUNCE        number of light bounces [default: 8]
  --direct DIRECT        number of direct rays to cast [default: 1]
  --indirect INDIRECT    number of indirect rays to cast [default: 1]
  --complete COMPLETE    number of samples-per-pixel at which to exit [default: +Inf]
  --thin                 treat transparent surfaces as having zero thickness
  --from FROM            camera position
  --to TO                camera target
  --focus FOCUS          camera focus (if other than 'to')
  --dist DIST            camera distance from target
  --polar POLAR          camera polar angle on target
  --longitude LONGITUDE
                         camera longitudinal angle on target
  --lens LENS            camera focal length in mm [default: 50]
  --fstop FSTOP          camera f-stop [default: 4]
  --expose EXPOSE        exposure multiplier [default: 1]
  --help, -h             display this help and exit
```

## More examples

I've gitignored the /fixtures directory to keep large binaries out of the repository.
You can download the Makefile fixtures from [Google Drive](https://drive.google.com/drive/folders/1hXQfQ9bZOIt8TvyoaUrRpELMxhKzrOCG?usp=sharing) for a library of models, materials, and HDR environments to play with.
