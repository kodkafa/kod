.PHONY: build clean

VERSION := $(shell git describe --tags --always --dirty)
COMMIT  := $(shell git rev-parse HEAD)
DATE    := $(shell date +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := -X kodkafa/internal/build.Version=$(VERSION) \
           -X kodkafa/internal/build.Commit=$(COMMIT) \
           -X kodkafa/internal/build.BuildDate=$(DATE)

build:
	go build -ldflags "$(LDFLAGS)" -o kod ./cmd/kod/main.go
	go build -ldflags "$(LDFLAGS)" -o kodkafa ./cmd/kodkafa/main.go

clean:
	rm -f kod kodkafa
