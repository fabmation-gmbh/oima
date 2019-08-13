BUILD_DATE ?= $(shell date --rfc-3339=seconds | sed -e 's/ /-/g')
COMMIT := $(shell git describe --dirty --always)
NUM_COMMITS := $(shell git rev-list --all --count)
VERSION ?= $(shell git describe --tags --dirty --always)
FULL_VERSION ?= "$(shell git describe --tags --dirty --always)-${NUM_COMMITS}-${COMMIT}"
LD_FLAGS := ""
DIST := $(CURDIR)/_dist

.PHONY: gen
gen:
	@cp -a internal/version.go.tpl internal/version.go
	@sed -i "s/commit unknown/$(COMMIT)/g" internal/version.go
	@sed -i "s/build date unknown/$(BUILD_DATE)/g" internal/version.go
	@sed -i "s/version unknown/$(VERSION)/g" internal/version.go

.PHONY: build
build: gen
	go build -o oima -ldflags $(LD_FLAGS)

.PHONY: clean
clean:
	rm -rf oima *.tar*

.PHONY: distclean
distclean: clean
	rm -rf $(DIST)

.PHONY: release
release: distclean gen
	mkdir -p $(DIST)
	@echo "+++++ Linux (amd64) +++++"
	GOOS=linux GOARCH=amd64 go build -o oima -ldflags $(LD_FLAGS)
	tar -zcvf $(DIST)/oima-$(VERSION)-linux-amd64.tar.gz oima README.md LICENSE CHANGELOG.md

	@echo "+++++ Windows (amd64) +++++"
	GOOS=windows GOARCH=amd64 go build -o oima.exe -ldflags $(LD_FLAGS)
	tar -zcvf $(DIST)/oima-$(VERSION)-windows-amd64.tar.gz oima.exe README.md LICENSE CHANGELOG.md

	@echo "+++++ Mac OSX +++++"
	GOOS=darwin GOARCH=amd64 go build -o oima -ldflags $(LD_FLAGS)
	tar -zcvf $(DIST)/oima-$(VERSION)-darwin-amd64.tar.gz oima README.md LICENSE CHANGELOG.md
