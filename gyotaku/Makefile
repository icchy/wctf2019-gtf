GOOS := linux
GOARCH := amd64

BUILD := $(shell git rev-parse --short HEAD)
LDFLAGS = -ldflags "-X=main.Build=$(BUILD)"


.PHONY: build
build-api:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/gyotaku $(LDFLAGS) -v ./src

.PHONY: windows
windows:
	GOOS=windows GOARCH=$(GOARCH) go build -o ./bin/gyotaku.exe $(LDFLAGS) -v ./src

.PHONY: clean
clean:
	rm -f bin/*
