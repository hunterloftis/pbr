package trace

import "os"

// Renderer renders the results of a trace to a file
type Renderer struct {
	Width  uint
	Height uint
}

func (r *Renderer) Write(file string) (n int, err error) {
	f, err := os.Create(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	n, err = f.WriteString("Hello\n")
	return n, err
}
