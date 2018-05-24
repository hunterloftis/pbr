package obj

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hunterloftis/pbr/material"
	"github.com/hunterloftis/pbr/rgb"
)

// http://paulbourke.net/dataformats/mtl/
type mtl struct {
	name  string
	kd    rgb.Energy // diffuse color
	tr    float64    // transmission
	ns    float64    // specular exponent (1 - roughness)
	ks    rgb.Energy // specular color
	ke    rgb.Energy // emissive color
	ni    float64    // refractive index (TODO: implement)
	pm    float64    // metal percent
	pc    float64    // clear-coat percent
	mapKd image.Image
}

func ReadMtl(filename string, thin bool) (mats map[string]*material.Map, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open mtl %v, %v", filename, err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	props := mtl{}
	mats = make(map[string]*material.Map)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := fields[0]
		args := fields[1:]
		switch key {
		case "newmtl":
			if len(props.name) > 0 {
				mats[props.name] = adapt(props)
			}
			props = mtl{name: args[0]}
		case "Kd":
			if props.kd, err = rgb.ParseEnergy(strings.Join(args, ",")); err != nil {
				return nil, err
			}
		case "Tr":
			if props.tr, err = strconv.ParseFloat(args[0], 64); err != nil {
				return nil, err
			}
		case "d":
			d, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return nil, err
			}
			props.tr = 1 - d
		case "Ns":
			if props.ns, err = strconv.ParseFloat(args[0], 64); err != nil {
				return nil, err
			}
		case "Ks":
			if props.ks, err = rgb.ParseEnergy(strings.Join(args, ",")); err != nil {
				return nil, err
			}
		case "Ke":
			if props.ke, err = rgb.ParseEnergy(strings.Join(args, ",")); err != nil {
				return nil, err
			}
		case "Ni":
			if props.ni, err = strconv.ParseFloat(args[0], 64); err != nil {
				return nil, err
			}
		case "Pm":
			if props.pm, err = strconv.ParseFloat(args[0], 64); err != nil {
				return nil, err
			}
		case "Pc":
			if props.pc, err = strconv.ParseFloat(args[0], 64); err != nil {
				return nil, err
			}
		case "map_Kd":
			rest := strings.TrimSpace(line[6:])
			texfile := filepath.Join(filepath.Dir(filename), rest)
			if props.mapKd, err = readTexture(texfile); err != nil {
				return nil, err
			}
		}
	}
	if len(props.name) > 0 {
		mats[props.name] = adapt(props)
	}
	return mats, nil
}

// https://github.com/AnalyticalGraphicsInc/obj2gltf#material-types
// http://exocortex.com/blog/extending_wavefront_mtl_to_support_pbr
// TODO: refractive index (ni) => .Fresnel
func adapt(props mtl) *material.Map {
	s := material.Sample{
		Color:        props.kd,
		Metalness:    props.pm,
		Roughness:    1 - (props.ns / 1000),
		Specularity:  0.04,
		Emission:     0,
		Transmission: props.tr,
	}
	if !props.ke.Zero() {
		rgb, scale := props.ke.Compressed(1)
		s.Color = rgb
		s.Emission = scale
	}
	m := material.MappedMaterial(s)
	// if props.ni > 1 {
	// 	// https://docs.blender.org/manual/en/dev/render/cycles/nodes/types/shaders/principled.html
	// 	f := math.Pow((props.ni-1)/(props.ni+1), 2)
	// 	d.Fresnel = rgb.Energy{f, f, f}
	// }
	return m
}

func readTexture(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open image %v, %v", filename, err)
	}
	defer f.Close()
	im, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return im, nil
}
