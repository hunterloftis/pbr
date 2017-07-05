go build -o bin/cubes ./cmd/cubes.go
bin/cubes -out out/render.png -heat out/heat.png
go tool pprof -text bin/cubes profile.pprof