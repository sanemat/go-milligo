VERSION = $(shell gobump show -r)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/sanemat/go-milligo.revision=$(CURRENT_REVISION)"
u := $(if $(update),-u)

.PHONY: test
test: download
	go test && \
	bash test.sh

.PHONY: download
download:
	go mod download && \
	go mod tidy

.PHONY: install-tools
install-tools: download
	$(MAKE) --directory=tools install-tools

.PHONY: goimports
goimports:
	goimports -w .

.PHONY: echo
echo:
	echo $(VERSION) $(BUILD_LDFLAGS)

.PHONY: build
build: download
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/milligo

.PHONY: install
install: download
	go install -ldflags=$(BUILD_LDFLAGS) ./cmd/milligo

.PHONY: crossbuild
crossbuild:
	goxz -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
      -os=linux,darwin,windows -d=./dist/v$(VERSION) ./cmd/*

.PHONY: upload
upload:
	ghr v$(VERSION) dist/v$(VERSION)

.PHONY: credits.txt
credits.txt:
	gocredits . > credits.txt

.PHONY: changelog
changelog:
	git-chglog -o changelog.md --next-tag v$(VERSION)
