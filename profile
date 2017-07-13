set -o errexit
set -x
mkdir -p out bin
go build -o bin/cubes cmd/cubes/cubes.go
bin/cubes -out out/render.png -heat out/heat.png -profile
go tool pprof -top bin/cubes profile.pprof
#go tool pprof -text -lines -total_delay -cum bin/cubes profile.pprof