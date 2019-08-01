BUILD_DATE ?= $(shell date --rfc-3339=seconds | sed -e 's/ /-/g')
COMMIT := $(shell git describe --dirty --always)
NUM_COMMITS := $(shell git rev-list --all --count)
VERSION ?= $(shell git describe --tags --dirty --always)
FULL_VERSION ?= "$(shell git describe --tags --dirty --always)-${NUM_COMMITS}-${COMMIT}"
LD_FLAGS := ""

.PHONY: gen
gen:
	@cp -a internal/version.go.tpl internal/version.go
	@sed -i "s/commit unknown/$(COMMIT)/g" internal/version.go
	@sed -i "s/build date unknown/$(BUILD_DATE)/g" internal/version.go
	@sed -i "s/version unknown/$(VERSION)/g" internal/version.go

.PHONY: build
build: gen
	go build -o oima -ldflags $(LD_FLAGS) .	