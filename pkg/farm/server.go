package farm

import (
	"fmt"
	"image/png"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/hunterloftis/pbr2/pkg/render"
)

// curl -H "Accept-Encoding: gzip" http://localhost:5000/scene > /dev/null

type Server struct {
	sample *render.Sample
}

func ListenAndServe(addr string, w, h int) error {
	s := NewServer(w, h)
	return s.ListenAndServe(addr)
}

func NewServer(w, h int) *Server {
	return &Server{
		sample: render.NewSample(w, h),
	}
}

func (s *Server) ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.getImage)
	mux.HandleFunc("/sample", s.postSample)
	return http.ListenAndServe(addr, handlers.CompressHandler(mux))
}

func (s *Server) getImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}
	if err := png.Encode(w, s.sample.Image()); err != nil {
		fmt.Fprintln(os.Stderr, "Error encoding png:", err)
		http.Error(w, "Unexpected error", 500)
	}
}

func (s *Server) postSample(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	err := s.sample.Read(r.Body)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Fprintln(w, "OK")
}
