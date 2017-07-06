set -o errexit
rm -f profile.pprof
go build -o bin/cubes ./cmd/cubes.go
bin/cubes -out out/render.png -heat out/heat.png -samples 8 -profile
go tool pprof -text bin/cubes profile.pprof