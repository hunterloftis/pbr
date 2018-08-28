.PHONY: hello redblue shapes toys

test:
	go test -v ./...
	go build ./cmd/pbr
	go build ./examples/shapes/shapes.go
	go build ./examples/hello/hello.go

hello:
	go run ./examples/hello/hello.go

shapes:
	go run ./examples/shapes/shapes.gom

fixtures:
	wget --no-check-certificate -r 'https://docs.google.com/uc?export=download&id=1mNizqg4uSHbCILBwf1-fbrJDu6dtvA7R' -O fixtures.zip
	unzip fixtures.zip
	rm -rf fixtures.zip __MACOSX

redblue:
	go run ./examples/redblue/redblue.go

sponza:
	go run ./examples/sponza/sponza.go

mario:
	go build ./cmd/pbr
	./pbr fixtures/models/mario/mario-sculpture.obj -from 100,100,400 -to 0,0,0 -env fixtures/envmaps/beach.hdr -rad 250 -floor 1.1 -width 1280 -height 720 -fstop 1.4 -out mario.png

lambo:
	go build ./cmd/pbr
	./pbr fixtures/models/lambo/lambo.obj -width 1280 -height 720 -env fixtures/envmaps/282.hdr -rad 2500 -to=-0.2,0.5,0.4 -from=-5,2,-5 -indirect -bounce 8

skull:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/skull.obj -scale 0.1,0.1,0.1 -floor 2 -material gold -env fixtures/envmaps/georgentor_4k.hdr -rad 100 -fstop 1.4 -from 2,0.2,1.75 -width 1280 -height 720

lucy:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/lucy.obj -scale 0.005,0.005,0.005 -material glass -env fixtures/envmaps/georgentor_4k.hdr -rad 100 -fstop 1.4 -width 1280 -height 720 -from 1,1.25,1 -to 0.1,1.25,0.1 -bounce 10

buddha:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/buddha.obj -width 1280 -height 720 -material gold -floor 20 -floorcolor 0,0,0 -env fixtures/envmaps/circus_maximus_1_4k.hdr -rad 200 -from 1.4,1.6,-1.8

falcon:
	go build ./cmd/pbr
	./pbr fixtures/models/falcon/millenium-falcon.obj -width 900 -height 450 -to=-86,-18,-2681 -from=500,300,-3400 -out falcon.png -env fixtures/envmaps/milkyway.hdr -rad 100 -sun 200,800,-3000 -sunsize 700 -lens 35

moses:
	go build ./cmd/pbr
	./pbr fixtures/models/moses/model.obj -out moses.png -scale 0.12,0.12,0.12 -width 1280 -height 720 -env fixtures/envmaps/georgentor_4k.hdr -rad 200 -fstop 1.4 -from 3,0,3 -to 0,0.5,0 -rotate 0,1,0

cesar:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/cesar.obj -width 500 -height 500

chair:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/chair.obj -width 480 -height 640 -from 40,300,-400

destroyer:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/destroyer.obj -width 1000 -height 400

legobricks:
	go build ./cmd/pbr
	./pbr fixtures/models/legobricks/LegoBricks3.obj -floor 2 -out legobricks.png

legoplane:
	go build ./cmd/pbr
	./pbr fixtures/models/legoplane/LEGO.Creator_Plane.obj -from 800,600,1300 -floor 10 -floorcolor 0.5,0.5,0.5 -floorrough 0.1 -env fixtures/envmaps/306.hdr -rad 1000 -width 1280 -height 720 -out legoplane.png

bowl:
	go build ./cmd/pbr
	./pbr fixtures/models/glassbowl/Glass\ Bowl\ with\ Cloth\ Towel.obj -from 6,4,6 -floor 2 -out bowl.png

glass:
	go build ./cmd/pbr
	./pbr fixtures/models/glass/glass-obj.obj -floor 1.5 -env fixtures/envmaps/ennis.hdr -from 840,120,600 -lens 80 -fstop 1.4 -focus 0.7 -out glass.png

toilet:
	go build ./cmd/pbr
	./pbr fixtures/models/toilet/Toilet.obj -width 320 -height 640 -from 0,200,150 -out toilet.png

gopher:
	go build ./cmd/pbr
	./pbr fixtures/models/gopher2/gopher.obj -floor 15

baccante1:
	go build ./cmd/pbr
	./pbr fixtures/models/baccante/baccante.obj -o baccante1.png -width 1280 -height 720 -floor 1 -rotate=-1.57,0,0 -scale 0.1,0.1,0.1 -to=0.5,0,-1.5 -from=-4.3,1,-2.6 -fstop 1.4 -env fixtures/envmaps/dresden_square_4k.hdr -rad 200

baccante2:
	go build ./cmd/pbr
	./pbr fixtures/models/baccante/baccante.obj -o baccante2.png -width 1280 -height 720 -floor 1 -rotate=-1.57,0,0 -scale 0.1,0.1,0.1 -to=0.5,0,-1.5 -from=-4.3,1,-2.6 -fstop 1.4 -env fixtures/envmaps/konzerthaus_4k.hdr -rad 200

baccante3:
	go build ./cmd/pbr
	./pbr fixtures/models/baccante/baccante.obj -o baccante3.png -width 1280 -height 720 -floor 1 -rotate=-1.57,0,0 -scale 0.1,0.1,0.1 -to=0.5,0,-1.5 -from=-4.3,1,-2.6 -fstop 1.4 -env fixtures/envmaps/venice_dawn_2_4k.hdr -rad 300

lion1:
	go build ./cmd/pbr
	./pbr fixtures/models/lion/lion.obj -out lion1.png -floor 1.1 -rotate=-1.57,0,0 -scale 0.2,0.2,0.2 -fstop 1.4 -env fixtures/envmaps/dresden_square_4k.hdr -rad 200 -width 1280 -height 720 -from 8,6,6

lion2:
	go build ./cmd/pbr
	./pbr fixtures/models/lion/lion.obj -out lion2.png -floor 1.1 -rotate=-1.57,0,0 -scale 0.2,0.2,0.2 -fstop 1.4 -env fixtures/envmaps/konzerthaus_4k.hdr -rad 300 -width 1280 -height 720 -from 8,6,6

lion3:
	go build ./cmd/pbr
	./pbr fixtures/models/lion/lion.obj -out lion3.png -floor 1.1 -rotate=-1.57,0,0 -scale 0.2,0.2,0.2 -fstop 1.4 -env fixtures/envmaps/venice_dawn_2_4k.hdr -rad 300 -width 1280 -height 720 -from 8,6,6

toys-fixtures:
	wget --no-check-certificate -r 'https://docs.google.com/uc?export=download&id=12YwRgYGilWMxtSek1uqF_ff8mfhpuYEB' -O toys.zip
	unzip toys.zip
	mkdir -p fixtures
	mv toys fixtures/toys
	rm -rf toys.zip __MACOSX
