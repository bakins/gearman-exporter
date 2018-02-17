NAME := gearman-exporter
PLATFORMS := linux/amd64 darwin/amd64
VERSION := $(shell git describe --tags --abbrev=0)

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

.PHONY: build
build: $(PLATFORMS)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -o ${NAME}.$(os).$(arch) ./cmd/${NAME}
