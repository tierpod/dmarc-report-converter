NAME       := dmarc-report-converter
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
build: lint $(NAME)

bin/$(NAME):
	go build -v $(LDFLAGS) -o $@ ./cmd/$(NAME)

.PHONY: clean
clean:
	rm -f bin/*
	rm -f ./pprof
	rm -rf ./tmp/dmarc-report-converter

.PHONY: install
install: $(INSTALLDIR)
	install -m 0755 bin/dmarc-report-converter $(INSTALLDIR)
	install -m 0600 config/config.dist.yaml $(INSTALLDIR)/config.dist.yaml
	cp -r templates $(INSTALLDIR)
	cp -r assets $(INSTALLDIR)
	cp -r install $(INSTALLDIR)

/etc/cron.daily/install/dmarc-report-converter.sh:
	install install/dmarc-report-converter.sh $@

$(INSTALLDIR):
	mkdir -p $@

.PHONY: uninstall
uninstall:
	rm -rf $(INSTALLDIR)
	rm -f /etc/cron.daily/dmarc-report-converter.sh

.PHONY: release
release: clean bin/$(NAME)
	mkdir -p release tmp
	make DESTDIR=./tmp install
	tar -cvzf release/$(NAME)_$(GIT_VER)_amd64.tar.gz --owner=0 --group=0 -C ./tmp $(NAME)
