NAME = docker2exe
OUTPUT = dist
VERSION = v0.1.1
SOURCES = $(wildcard *.go)
TARGETS = darwin/amd64 linux/amd64 windows/amd64

OS = $(shell go env GOOS)
ARCH = $(shell go env GOARCH)

os = $(word 1, $(subst /, ,$@))
arch = $(word 2, $(subst /, ,$@))

.PHONY: all
all: $(TARGETS)

.PHONY: test
test: all
	dist/docker2exe-$(OS)-$(ARCH) --name alpine --image alpine
	dist/alpine-$(OS)-$(ARCH) echo OK
	dist/docker2exe-$(OS)-$(ARCH) --name alpine-embed --image alpine --embed
	dist/alpine-embed-$(OS)-$(ARCH) echo OK

.PHONY: release
release: all
	hub release create $(VERSION) \
		-a dist/docker2exe-linux-amd64 \
		-a dist/docker2exe-darwin-amd64 \
		-a dist/docker2exe-windows-amd64

$(OUTPUT):
	mkdir $(OUTPUT)

$(TARGETS): $(SOURCES) $(OUTPUT)
	GOOS=$(os) GOARCH=$(arch) go build -o "$(OUTPUT)/$(NAME)-$(os)-$(arch)"
