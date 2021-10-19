export tag=0.7
CGO_ENABLED := 0
# GOOS := darwin
GOOS := linux
GOARCH := amd64
PORT := 9000
build:
	echo "building go-httpserver binary"
	echo $(GOOS)
	mkdir -p bin/amd64
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(GOARCH) .

release:
	echo "building httpserver container"
	docker build -t httpserver:${tag} .

run: 
	echo "running httpserver"
	docker run --rm -it -p $(PORT):$(PORT) httpserver:${tag} 

push: release
	echo "pushing httpserver"
	docker tag httpserver:${tag} yueyongqiang/httpserver:${tag}
	docker push yueyongqiang/httpserver:${tag}