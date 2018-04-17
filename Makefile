BINARIES  := bin/dmarc-report-converter

VERSION := $(shell git describe --tags)

LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: lint
lint:
	find ./cmd ./pkg -type f -name '*.go' | xargs gofmt -l -e
	go vet ./cmd/... ./pkg/...
	$(GOPATH)/bin/golint ./cmd/... ./pkg/...
	#go test ./cmd/... ./pkg/...

.PHONY: build
build: lint $(BINARIES)

$(BINARIES):
	go build -v $(LDFLAGS) -o $@ cmd/$(notdir $@)/*.go

.PHONY: clean
clean:
	rm -f bin/*
	rm -f install/*.retry
	rm -f ./pprof

.PHONY: doc
doc:
	godoc -http :6060
