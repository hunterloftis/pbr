package render

type Scene struct {
	Camera  Camera
	Env     Environment
	Surface Surface
}

func NewScene(c Camera, s Surface, e Environment) *Scene {
	return &Scene{
		Env:     e,
		Surface: s,
		Camera:  c,
	}
}

func (s *Scene) Render(width, height, bounce int, direct bool) *Frame {
	f := NewFrame(s, width, height, bounce, direct)
	f.Start()
	return f
}
