export tag=0.3
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

release: build
	echo "building httpserver container"
	docker build -t httpserver:${tag} .

run: release
	echo "pushing httpserver"
	docker run -it -p $(PORT):$(PORT) httpserver:${tag} --rm