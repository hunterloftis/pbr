set -o errexit
set -x
mkdir -p out bin
go build -o bin/cubes cmd/cubes/cubes.go
bin/cubes -out out/render.png -heat out/heat.png -profile block
go tool pprof -top -lines bin/cubes profile.pprof