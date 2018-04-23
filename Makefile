BINARIES := bin/dmarc-report-converter
DESTDIR  := /opt
VERSION  := $(shell git describe --tags)
LDFLAGS  := -ldflags "-X main.version=$(VERSION)"

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

.PHONY: install
install: $(DESTDIR)/dmarc-report-converter
	install -m 0755 bin/dmarc-report-converter $(DESTDIR)/dmarc-report-converter
	install -m 0600 config/config.dist.yaml $(DESTDIR)/dmarc-report-converter/config.dist.yaml
	install install/dmarc-report-converter.cron /etc/cron.daily/

$(DESTDIR)/dmarc-report-converter:
	mkdir -p $@

.PHONY: uninstall
uninstall:
	rm -rf $(DESTDIR)/dmarc-report-converter
	rm -f /etc/cron.daily/dmarc-report-converter.cron
