.PHONY: fixtures

clean:
	rm -f *.png *.pdf *.zip

count:
	find . -name '*.go' -not -path "./vendor/*" | xargs wc -l

fixtures:
	curl -L -o fixtures.zip https://www.dropbox.com/sh/ik2vfz1qhtsupgt/AADeLgXNSrcjkhqbY64ng5bRa?dl=1
	unzip -n fixtures.zip -d fixtures -x /

doc:
	godoc -http=":5000"

pbr:
	go build ./cmd/pbr

help: pbr
	./pbr --help

hello:
	go run ./examples/hello/hello.go
	open hello.png

adaptive: pbr
	./pbr fixtures/models/falcon.obj -o nonadaptive.png -dist 480 -lat 0.25 -lon=-1 -target=-86,-55,-2770 -focus=-86,-18,-2682 -heat nonadaptive-heat.png -width 888 -height 300 -branch 1 -adapt 0 -time 600
	./pbr fixtures/models/falcon.obj -o adaptive.png -dist 480 -lat 0.25 -lon=-1 -target=-86,-55,-2770 -focus=-86,-18,-2682 -heat adaptive-heat.png -width 888 -height 300 -time 600
	open adaptive.png adaptive-heat.png nonadaptive.png nonadaptive-heat.png

shapes:
	go run ./examples/shapes/shapes.go
	open shapes.png
	
destroyer: pbr
	./pbr fixtures/models/destroyer.obj -dist 12000 -lon 0.4 -width 1200 -height 500 -complete 8
	open destroyer.png
	
skull: pbr
	./pbr fixtures/models/skull.obj -complete 16
	open skull.png

chair: pbr
	./pbr fixtures/models/chair.obj -lens 150 -fstop 2.8 -floor --complete 10
	open chair.png

lambo: pbr
	./pbr fixtures/models/lambo3/lambo.obj -heat lambo-heat.png -lon 4 -lat 0.1 -env fixtures/images/river.hdr -rad 2500 -lens 50 -fstop 1.4 -target=-0.2,0.5,0.4 -dist 6.6 -focus=-1,0.67,-0.56 -direct 0 -width 1920 -height 1080 -complete 256
	open lambo.png

profile: pbr
	./pbr fixtures/models/lambo3/lambo.obj -heat lambo-heat.png -lon 4 -lat 0.1 -env fixtures/images/river.hdr -rad 2500 -lens 50 -fstop 1.4 -target=-0.2,0.5,0.4 -dist 6.6 -focus=-1,0.67,-0.56 -direct 0 -width 1920 -height 1080 -complete 1 -profile
	go tool pprof --pdf ./pbr ./cpu.pprof > cpu.pdf
	open cpu.pdf

ibl: pbr
	./pbr fixtures/models/mario/mario-sculpture.obj -o mario-249.png -lon 1 -lat 0.1 -floor -env fixtures/images/249.hdr -rad 600 -width 888 -height 500 -complete 50
	./pbr fixtures/models/mario/mario-sculpture.obj -o mario-beach.png -lon 1 -lat 0.1 -floor -env fixtures/images/beach.hdr -rad 350 -width 888 -height 500 -complete 50
	./pbr fixtures/models/mario/mario-sculpture.obj -o mario-misty.png -lon 1 -lat 0.1 -floor -env fixtures/images/misty.hdr -rad 300 -width 888 -height 500 -complete 70
	./pbr fixtures/models/mario/mario-sculpture.obj -o mario-lobe.png -lon 1 -lat 0.1 -floor -env fixtures/images/lobe.hdr -rad 400 -width 888 -height 500 -complete 60
	./pbr fixtures/models/mario/mario-sculpture.obj -o mario-ambient.png -lon 1 -lat 0.1 -floor -width 888 -height 500 -complete 50
	open mario-249.png mario-beach.png mario-misty.png mario-lobe.png mario-ambient.png
