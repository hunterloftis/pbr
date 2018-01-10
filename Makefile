test:
	go build ./cmd/pbr
	./pbr fixtures/models/lucy.obj
	open lucy.png

chair:
	go build ./cmd/pbr
	./pbr fixtures/models/chair.obj chair.png -heat heat.png -noise noise.png -from 400,400,0 -to "10,75,-600" -sky 200,250,300 -lens 150 -fstop 2.8 -profile --complete 4
	go tool pprof --pdf ./pbr ./cpu.pprof > profile.pdf && open profile.pdf

skull:
	go build ./cmd/pbr
	./pbr fixtures/models/skull.obj skull.png -from "9,2.46,10" -to "0,2.46,-0.59" -env fixtures/images/glacier.hdr -lens 35 -expose 2 -complete 100
	open skull.png

lucy:
	go build ./cmd/pbr
	./pbr fixtures/models/lucy.obj lucy.png -from "800,400,500" -to "15,165,-4" -env fixtures/images/uffizi.hdr -lens 70 -fstop 0.01 -width 600 -height 750 -expose 2 -noise noise.png -heat heat.png -complete 1000
	open lucy.png heat.png noise.png

destroyer:
	go build ./cmd/pbr
	./pbr fixtures/models/destroyer.obj -dist 12000 -polar 0 -width 1200 -height 500 -complete 10 -profile
	open destroyer.png
	go tool pprof --pdf ./pbr ./cpu.pprof > profile.pdf && open profile.pdf

bmw:
	go build ./cmd/pbr
	./pbr fixtures/models/bmw/BMW850.obj -polar=-1 -longitude=0.3 -width 1200 -height 500 -lens 35 -dist 80 -complete 100
	open BMW850.png
