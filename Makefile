BINARIES   := bin/dmarc-report-converter
DESTDIR    := /opt
INSTALLDIR := $(DESTDIR)/dmarc-report-converter

GIT_VER    := $(shell git describe --abbrev=7 --always --tags)-$(shell git rev-parse --abbrev-ref HEAD)-$(shell date +%Y%m%d)
LDFLAGS    := -ldflags "-X main.version=$(GIT_VER)"

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
	rm -f release/*.tar.gz
	rm -f install/*.retry
	rm -f ./pprof

.PHONY: doc
doc:
	godoc -http :6060

.PHONY: install install
install: $(INSTALLDIR) $(INSTALLDIR)/templates
	install -m 0755 bin/dmarc-report-converter $(INSTALLDIR)
	install -m 0600 config/config.dist.yaml $(INSTALLDIR)/config.dist.yaml
	cp -r templates $(INSTALLDIR)

/etc/cron.daily/install/dmarc-report-converter.cron:
	install install/dmarc-report-converter.cron $@

$(INSTALLDIR) $(INSTALLDIR)/templates:
	mkdir -p $@

.PHONY: uninstall
uninstall:
	rm -rf $(DESTDIR)/dmarc-report-converter
	rm -f /etc/cron.daily/dmarc-report-converter.cron

release/dmarc-report-converter_linux_amd64.tar.gz: $(BINARIES)
	mkdir -p release
	make DESTDIR=./tmp install
	tar -cvzf $@ --owner=0 --group=0 -C ./tmp dmarc-report-converter
