.PHONY: fixtures

adaptive:
	go build ./cmd/pbr
	./pbr fixtures/models/falcon.obj -o nonadaptive.png -dist 480 -lat 0.25 -lon=-1 -target=-86,-55,-2770 -focus=-86,-18,-2682 -heat nonadaptive-heat.png -width 888 -height 300 -branch 1 -adapt 0 -time 600
	./pbr fixtures/models/falcon.obj -o adaptive.png -dist 480 -lat 0.25 -lon=-1 -target=-86,-55,-2770 -focus=-86,-18,-2682 -heat adaptive-heat.png -width 888 -height 300 -time 600
	open adaptive.png adaptive-heat.png nonadaptive.png nonadaptive-heat.png

shapes:
	go run ./examples/shapes.go
	open shapes.png
	
help:
	go build ./cmd/pbr
	./pbr --help
	
hello:
	go run ./examples/hello.go
	open hello.png
	
destroyer:
	go build ./cmd/pbr
	./pbr fixtures/models/destroyer.obj -dist 12000 -lon 0.4 -width 1200 -height 500 -complete 8
	open destroyer.png

falcon:
	go build ./cmd/pbr
	./pbr fixtures/models/falcon.obj -dist 900 -lat 0.3 -lon=-1 -heat falcon-heat.png -complete 8
	open falcon.png falcon-heat.png

count:
	find . -name '*.go' -not -path "./vendor/*" | xargs wc -l

fixtures:
	@echo "download https://drive.google.com/drive/folders/1hXQfQ9bZOIt8TvyoaUrRpELMxhKzrOCG?usp=sharing into ./fixtures"

doc:
	godoc -http=":5000"
	
skull:
	go build ./cmd/pbr
	./pbr fixtures/models/skull.obj -complete 100
	open skull.png

chair:
	go build ./cmd/pbr
	./pbr fixtures/models/chair.obj -lens 150 -fstop 2.8 -floor --complete 10
	open chair.png

bmw:
	go build ./cmd/pbr
	./pbr fixtures/models/bmw/BMW850.obj -lon=-1.2 -lat=0.3 -width 1200 -height 500 -lens 35 -dist 85 -thin -floor -target 7.33,11,-99.21 -env fixtures/images/pisa.hdr -rad 1600 -complete 500
	open BMW850.png

bmw2:
	go build ./cmd/pbr
	./pbr fixtures/models/bmw/BMW850.obj -lon=-2.1 -lat=0.08 -width 1200 -height 500 -lens 65 -dist 120 -thin -floor -target 7.33,9,-99.21 -focus 7.33,8,-140 -env fixtures/images/glacier.hdr -rad 400 -fstop 0.1 -complete 32
	open BMW850.png

lambo:
	go build ./cmd/pbr
	./pbr fixtures/models/lambo2/lambo.obj -heat lambo-heat.png -floor -lon 4 -lat 0.1 -env fixtures/images/306.hdr -rad 700 -lens 50 -fstop 1.4 -target=-0.2,0.5,0.4 -dist 6.5 -focus=-1,0.67,-0.56 -direct 0 -width 1920 -height 1080 -complete 512
	open lambo.png

profile:
	go build ./cmd/pbr
	./pbr fixtures/models/lambo2/lambo.obj -floor -lon 3.6 -lat 0.1 -env fixtures/images/293.hdr -rad 450 -lens 60 -fstop 1.4 -target=-0.1,0.5,0.1 -dist 7.5 -focus=-2.2658,0.5542,-0.7 -direct 0 -width 960 -height 540 -profile -complete 1
	go tool pprof --pdf ./pbr ./cpu.pprof > cpu.pdf
	open cpu.pdf

ibl:
	go build ./cmd/pbr
	./pbr fixtures/models/mario/mario-sculpture.obj -o mario-249.png -lon 1 -lat 0.1 -floor -env fixtures/images/249.hdr -rad 600 -width 888 -height 500 -complete 64
	./pbr fixtures/models/mario/mario-sculpture.obj -o mario-beach.png -lon 1 -lat 0.1 -floor -env fixtures/images/beach.hdr -rad 300 -width 888 -height 500 -complete 64
	./pbr fixtures/models/mario/mario-sculpture.obj -o mario-misty.png -lon 1 -lat 0.1 -floor -env fixtures/images/misty.hdr -rad 300 -width 888 -height 500 -complete 64
	./pbr fixtures/models/mario/mario-sculpture.obj -o mario-lobe.png -lon 1 -lat 0.1 -floor -env fixtures/images/lobe.hdr -rad 300 -width 888 -height 500 -complete 64
	open mario-249.png mario-beach.png mario-misty.png mario-lobe.png
