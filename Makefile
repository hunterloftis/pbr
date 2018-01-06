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
	./pbr fixtures/models/lucy.obj lucy.png -from "900,300,600" -to "15,165,-4" -env fixtures/images/uffizi.hdr -lens 90 -fstop 0.01 -width 600 -height 750 -expose 2 -noise noise.png -heat heat.png -adapt 5 -complete 500
	open lucy.png heat.png noise.png
