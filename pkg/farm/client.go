package farm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hunterloftis/pbr2/pkg/render"
)

func Render(scene *render.Scene, url string, w, h, depth int, direct bool) error {
	frame := scene.Render(w, h, depth, direct)
	defer frame.Stop()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	fmt.Printf("\nRendering to %v", url)

	for frame.Active() {
		<-ticker.C
		if sample, n := frame.Sample(); n > 0 {
			fmt.Print(".")
			buf, err := sample.Buffer()
			if err != nil {
				fmt.Println("\nError:", err)
				continue
			}
			resp, err := http.Post(url, "application/octet-stream", buf) // TODO: gzip
			if resp.StatusCode != 200 || err != nil {
				fmt.Println("\nStatus:", resp.Status, "Error:", err)
				continue
			}
			frame.Clear()
		}
	}
	return nil
}
