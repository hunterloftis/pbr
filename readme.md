# A Physically-Based Renderer in Go

Learning Go by writing a path tracer.

![Render](https://user-images.githubusercontent.com/364501/27873088-e36c621c-6178-11e7-94c8-19171e05dc50.png)

- Unbiased Monte-Carlo integration
- Adaptive sampling
- Russian roulette (early path termination)
- Parametric shapes (spheres, cubes)
- Physically-based materials
  - Fresnel reflection, transmission, absorption, diffusion
  - Color, refractive indices, gloss, transparency, separate fresnel channels, metals
- Arbitrary light sources ('everything is a light')
- Environment maps
- Physically-based cameras
  - Sensor, aperture, focus, depth-of-field
- Anti-aliasing

## Try it

```
$ go get -u github.com/hunterloftis/pbr/pbr
$ cd $GOPATH/src/github.com/hunterloftis/pbr
$ ./run
```

By default, the renderer runs until it receives a signal (like Ctrl + C)