package obj

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hunterloftis/pbr/surface"
)

// ReadFile reads the 3D geometry data from a Wavefront .obj file.
// It automatically reads material data from any referenced .mtl files.
// Missing material data is not an error; missing .mtl files are skipped.
func ReadFile(filename string, thin bool) ([]surface.Surface, error) {
	s := make([]surface.Surface, 0)
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open scene %v, %v", filename, err)
	}
	defer f.Close()
	scanner := NewScanner(f)
	for scanner.Scan() {
		switch n := scanner.Next().(type) {
		case MatName:
			mfile := filepath.Join(filepath.Dir(filename), string(n))
			mats, err := ReadMtl(mfile, thin)
			if err != nil {
				return nil, err
			}
			scanner.AddMaterials(mats)
		case TexName:
		case surface.Surface:
			s = append(s, n)
		}
	}
	return s, scanner.Err()
}
