all: certdump image

certdump: certdump.go
	go build  -a -tags "netgo static_build" -installsuffix netgo ./certdump.go


image: certdump Dockerfile
	docker build -t dhiltgen/certdump:latest .

push:
	docker push dhiltgen/certdump:latest


PHONY: image
