test:
	go build ./cmd/pbr
	./pbr fixtures/models/lucy.obj -complete 2
	open lucy.png

chair:
	go build ./cmd/pbr
	./pbr fixtures/models/chair.obj -from 400,400,0 -to 10,75,-600 -sky 200,250,300 -lens 150 -fstop 2.8 -floor --complete 10
	open chair.png

skull:
	go build ./cmd/pbr
	./pbr fixtures/models/skull.obj -from 11,4,12 -to=-0.57,3.2,-1.69 -env fixtures/images/glacier.hdr -lens 50 -fstop 0.3 -expose 2 -complete 50
	open skull.png

lucy:
	go build ./cmd/pbr
	./pbr fixtures/models/lucy.obj lucy.png -from "800,400,500" -to "15,165,-4" -env fixtures/images/uffizi.hdr -lens 70 -fstop 0.01 -width 600 -height 750 -expose 2 -noise noise.png -heat heat.png -complete 1000
	open lucy.png heat.png noise.png

destroyer:
	go build ./cmd/pbr
	./pbr fixtures/models/destroyer.obj -dist 12000 -longitude 0 -width 1200 -height 500 -complete 10 -profile
	open destroyer.png
	go tool pprof --pdf ./pbr ./cpu.pprof > profile.pdf && open profile.pdf

bmw:
	go build ./cmd/pbr
	./pbr fixtures/models/bmw/BMW850.obj -polar=-1.2 -longitude=0.3 -width 1200 -height 500 -lens 35 -dist 85 -thin -floor -to 7.33,11,-99.21 -env fixtures/images/pisa.hdr -rad 300 -complete 500
	open BMW850.png

bmw2:
	go build ./cmd/pbr
	./pbr fixtures/models/bmw/BMW850.obj -polar=-2.1 -longitude=0.08 -width 1200 -height 500 -lens 65 -dist 120 -thin -floor -to 7.33,9,-99.21 -focus 7.33,8,-140 -env fixtures/images/glacier.hdr -rad 400 -fstop 0.1 -complete 32
	open BMW850.png

lambo:
	go build ./cmd/pbr
	./pbr fixtures/models/lambo2/lamborghini-aventador-pbribl.obj -floor -out lambo.png -polar 3.6 -longitude 0.1 -env fixtures/images/293.hdr -rad 400 -width 1152 -height 648 -lens 60 -fstop 1 -to=-0.1,0.25,0 -dist 7.5 -bounce 7 -complete 5000
	open lambo.png
