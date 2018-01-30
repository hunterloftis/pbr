# pbr: a Physically-Based 3D Renderer in Go

Package pbr implements Physically-Based Rendering with a Monte Carlo path tracer.
[[ Documentation ]](https://godoc.org/github.com/hunterloftis/pbr)
[[ Github ]](https://github.com/hunterloftis/pbr)

![Lambo Render](https://user-images.githubusercontent.com/364501/35541108-27ce146c-0527-11e8-83e2-1af8aee2c0e2.png)

```
$ make fixtures
$ make lambo
```

This is an unbiased forward path-tracer written in Go and inspired by Disney's [Hyperion video](https://www.disneyanimation.com/technology/innovations/hyperion). It traces light rays from
physically-based cameras to surfaces with realistic materials and simulates how the rays should bounce around the scene.
It has a simple API for creating scenes in code and a CLI for rendering photorealistic images from the command line.

## Try it

Install:

```
$ go get github.com/hunterloftis/pbr
$ cd $GOPATH/src/github.com/hunterloftis/pbr
$ dep ensure
```

Run:

```
$ cd $GOPATH/src/github.com/hunterloftis/pbr
$ go build ./cmd/pbr
$ pbr fixtures/models/falcon.obj -lat 0.5 -lon 0.5 -complete 5
$ open falcon.png
```

## Features

### Shapes & Transforms

![shapes](https://user-images.githubusercontent.com/364501/35257181-c771dd1c-ffc5-11e7-96d9-0a576a886b3c.png)

```
$ make shapes
```

The renderer supports spheres, cubes, and triangles that can be moved, scaled, and rotated.

### Image-based lighting

![ibl](https://user-images.githubusercontent.com/364501/35474229-60806258-0359-11e8-8af8-a67901ae9720.png)

```
$ make fixtures
$ make ibl
```

The renderer can use high dynamic range (HDR) panoramic images as complex, omnidirectional light sources.
This allows highly detailed real-world lighting to illuminate the scene's surfaces for greater realism and visual interest.

All the above images were rendered with an identical model and identical materials;
the only difference between them is the image used for lighting.

### Sampling & Branching

![falcon adaptive](https://user-images.githubusercontent.com/364501/35202761-753e2d44-fef2-11e7-8d55-4893eb860144.png)
![falcon nonadaptive](https://user-images.githubusercontent.com/364501/35202760-752b55ca-fef2-11e7-8181-e77e137c1668.png)

```
$ make adaptive
```

Adaptive sampling devotes more time to sampling noisy areas than already-resolved ones.
Branched tracing splits primary rays into multiple branches to better sample the most important (first) bounce of each path.
Both of these techniques allow the renderer to spend its Ray-Scene intersection budget more effectively.

Both closeups of the Millennium Falcon were rendered in 10 minutes.
The top image used naive sampling while the bottom used the default adaptive and branching settings.

### Wavefront .obj files

### Physically-based materials

### Arbitrary light sources

### Direct lighting

### Thin surfaces

### Physically-based cameras

### Supersampled anti-aliasing

### K-D Tree acceleration

### Summary

- Geometry:
  - [Parametric shapes (spheres, cubes, triangles)](#shapes--transforms)
  - [Transformation matrices (translate, rotate, scale)](#shapes--transforms)
  - Wavefront .obj files (meshes) and .mtl files (materials)
- Materials:
  - [Physically-based materials](https://www.marmoset.co/posts/basic-theory-of-physically-based-rendering/)
  - [PBR extensions for .mtl files](http://exocortex.com/blog/extending_wavefront_mtl_to_support_pbr)
  - Reflection, transmission, absorption, diffusion
  - Color, refractive index, roughness, transparency, fresnel channels, metalness, clear coat
  - Texture maps (.jpg, .png)
- Lighting:
  - [Direct lighting]()
  - Arbitrary light sources ('everything is a light')
  - [Image-based lighting](#image-based-lighting)
- Cameras:
  - Physically-based cameras
  - Sensor, aperture, focal length, focus, depth-of-field
- Quality and speed:
  - [Adaptive sampling](#sampling--branching)
  - [Branched tracing](#sampling--branching)
  - [Russian roulette](https://computergraphics.stackexchange.com/questions/2316/is-russian-roulette-really-the-answer)
  - [K-D Tree acceleration](http://slideplayer.com/slide/7653218/)
  - [Supersampled anti-aliasing](https://en.wikipedia.org/wiki/Supersampling)
- Interface:
  - 100% Go with no system dependencies
  - Sequential API, concurrent execution
  - CLI

## API

See [GoDoc](https://godoc.org/github.com/hunterloftis/pbr)

## CLI

```
Usage: pbr SCENE [options]

Positional arguments:
  SCENE                  input scene .obj

Options:
  --verbose              verbose output with scene information
  --info, -i             output scene information and exit
  --out OUT, -o OUT      output render .png
  --heat HEAT            output heatmap as .png
  --noise NOISE          output noisemap as .png
  --profile              record performance into profile.pprof
  --width WIDTH          rendering width in pixels [default: 800]
  --height HEIGHT        rendering height in pixels [default: 450]
  --target TARGET        camera target point
  --focus FOCUS          camera focus point (if other than 'target')
  --dist DIST            camera distance from target
  --lat LAT              camera polar angle on target
  --lon LON              camera longitudinal angle on target
  --lens LENS            camera focal length in mm [default: 50]
  --fstop FSTOP          camera f-stop [default: 4]
  --expose EXPOSE        exposure multiplier [default: 1]
  --ambient AMBIENT      the ambient light color [default: &{500 500 500}]
  --env ENV, -e ENV      environment as a panoramic hdr radiosity map (.hdr file)
  --rad RAD              exposure of the hdr (radiosity) environment map [default: 100]
  --floor                create a floor underneath the scene
  --adapt ADAPT          adaptive sampling multiplier [default: 8]
  --bounce BOUNCE, -d BOUNCE
                         number of light bounces (depth) [default: 8]
  --direct DIRECT, -d DIRECT
                         maximum number of direct rays to cast [default: 1]
  --branch BRANCH, -b BRANCH
                         maximum number of branches on first hit [default: 32]
  --complete COMPLETE, -c COMPLETE
                         number of samples-per-pixel at which to exit [default: +Inf]
  --thin                 treat transparent surfaces as having zero thickness
  --help, -h             display this help and exit
  --version              display version and exit
```
