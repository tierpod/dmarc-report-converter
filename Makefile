NAME       := dmarc-report-converter
DESTDIR    := /opt
INSTALLDIR := $(DESTDIR)/dmarc-report-converter
ARCH       := $(shell arch)

ifeq ($(GITHUB_REF),)
GIT_VER    := $(shell git describe --abbrev=7 --always --tags)-$(shell git rev-parse --abbrev-ref HEAD)-$(shell date +%Y%m%d)
else
GIT_VER    := $(shell basename $(GITHUB_REF))-$(shell date +%Y%m%d)
endif
LDFLAGS    := -ldflags "-X main.version=$(GIT_VER)"

CGO_ENABLED := 0

.PHONY: test
test:
	find ./cmd ./pkg -type f -name '*.go' | xargs gofmt -l -e
	go vet -mod=vendor ./cmd/... ./pkg/...
	$(shell go env GOPATH)/bin/staticcheck ./cmd/... ./pkg/...
	go test -mod=vendor ./cmd/... ./pkg/...

.PHONY: build
build: test bin/$(NAME)

bin/$(NAME):
	CGO_ENABLED=$(CGO_ENABLED) go build -mod=vendor -v $(LDFLAGS) -o $@ ./cmd/$(NAME)

.PHONY: clean
clean:
	rm -f bin/*
	rm -f ./pprof
	rm -rf ./tmp/dmarc-report-converter

.PHONY: install
install: $(INSTALLDIR) bin/$(NAME)
	install -m 0755 bin/$(NAME) $(INSTALLDIR)
	install -m 0600 config/config.dist.yaml $(INSTALLDIR)/config.dist.yaml
	cp -r assets $(INSTALLDIR)
	cp -r install $(INSTALLDIR)

$(INSTALLDIR) dist tmp:
	mkdir -p $@

.PHONY: release
release: clean dist
	make DESTDIR=./tmp install
	tar -cvzf dist/$(NAME)_$(GIT_VER)_$(ARCH).tar.gz --owner=0 --group=0 -C ./tmp $(NAME)
