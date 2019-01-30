# Makefile for project binnit using go mod.


#V := 1 # When V is set, print commands and build progress.

all: release

build:
	$Q for arch in 386 amd64 arm64 ; do \
	CGO_ENABLED=0 GOOS=linux GOARCH=$$arch go build $(VERSION_FLAGS) -a -installsuffix cgo -o ./bin/binnit-linux-$$arch $(IMPORT_PATH); \
	done

release: clean build
	$Q mkdir -p ./release
	$Q for arch in 386 amd64 arm64 ;  do \
	rm -rf ./tmp/binnit ; \
	mkdir -p ./tmp/binnit ; \
	cp -a ./bin/binnit-linux-$$arch ./tmp/binnit/binnit ; \
	cp -a ./tpl ./tmp/binnit ; \
	cp -a ./static ./tmp/binnit ; \
	mkdir -p ./tmp/binnit/paste ; \
	mkdir -p ./tmp/binnit/log ; \
	mkdir -p ./tmp/binnit/conf ; \
	cp -a ./binnit.cfg ./tmp/binnit/conf ; \
	tar -cz -C ./tmp -f ./release/binnit-$$arch-$(VERSION).tar.gz ./binnit ; \
	done

##### =====> Utility targets <===== #####

.PHONY: clean gen

clean:
	$Q rm -rf bin
	$Q rm -rf $(wildcard ./release/*-$(VERSION)*)

gen:
	@echo "Running go generate"
	$Q go generate
	@echo "Done!"

.PHONY: all build release

##### =====> Internals <===== #####

Q                := $(if $V,,@)
IMPORT_PATH      := $(shell awk -F" " '$$1=="module"{print $$2;exit;}' go.mod)
VERSION          := $(shell git describe --tags --always --dirty="-dev")
DATE             := $(shell date -u '+%Y-%m-%d-%H%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'
