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

func ReadMtl(filename string, thin bool) (mats []*material.Map, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open mtl %v, %v", filename, err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	props := mtl{}
	mats = make([]*material.Map, 0)
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
				mats = append(mats, adapt(props, thin))
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
		mats = append(mats, adapt(props, thin))
	}
	return mats, nil
}

// https://github.com/AnalyticalGraphicsInc/obj2gltf#material-types
// http://exocortex.com/blog/extending_wavefront_mtl_to_support_pbr
// TODO: refractive index (ni) => .Fresnel
func adapt(props mtl, thin bool) *material.Map {
	d := material.MaterialDesc{
		Name:     props.name,
		Color:    props.kd,
		Transmit: props.tr,
		Rough:    1 - (props.ns / 1000),
		Light:    props.ke,
		Metal:    props.pm,
		Thin:     thin,
		Coat:     props.pc,
		Texture:  props.mapKd,
	}
	if props.tr > 0 {
		if d.Thin {
			d.Transmit = props.tr
		} else {
			d.Transmit = 1
			d.Color = d.Color.Amplified(props.tr)
		}
		d.Fresnel = rgb.Energy{0.042, 0.042, 0.042} // Glass
	} else {
		d.Fresnel = rgb.Energy{0.02, 0.02, 0.02}.Blend(props.ks, d.Metal)
	}
	return material.New(d)
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
