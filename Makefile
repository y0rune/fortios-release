# Source from https://github.com/krisnova/go-nova/blob/main/Makefile

all: compile

version     ?=  @(git describe --tags --abbrev=0)
org         ?=  y0rune
target		?=  fortios-release
authorname  ?=  Marcin Woźniak
authoremail ?=  y0rune@aol.com
license     ?=  MIT
year        ?=  2022
copyright   ?=  Copyright (c) $(year)
gofile      ?= main.go

compile:
	@echo "Compiling..."
	go build -ldflags "\
	-X 'github.com/$(org)/$(target).Version=$(version)' \
	-X 'github.com/$(org)/$(target).AuthorName=$(authorname)' \
	-X 'github.com/$(org)/$(target).AuthorEmail=$(authoremail)' \
	-X 'github.com/$(org)/$(target).Copyright=$(copyright)' \
	-X 'github.com/$(org)/$(target).License=$(license)' \
	-X 'github.com/$(org)/$(target).Name=$(target)'" \
	-o $(target) $(gofile)

# install:
# 	@echo "Installing..."
# 	sudo cp $(target) /usr/bin/$(target)

build: clean compile

clean:
	@echo "Cleaning..."
	rm -rvf release/*
	rm -rvf $(target)

# test: clean compile install
# test: clean compile
# 	@echo "Testing..."
# 	go test -v .

.PHONY: release
release:
	mkdir -p release
	GOOS="linux" GOARCH="amd64" go build -ldflags "-X 'github.com/$(org)/$(target).Version=$(version)'" -o release/$(target)-linux-amd64 $(gofile)
	GOOS="linux" GOARCH="arm" go build -ldflags "-X 'github.com/$(org)/$(target).Version=$(version)'" -o release/$(target)-linux-arm $(gofile)
	GOOS="linux" GOARCH="arm64" go build -ldflags "-X 'github.com/$(org)/$(target).Version=$(version)'" -o release/$(target)-linux-arm64 $(gofile)
	GOOS="darwin" GOARCH="amd64" go build -ldflags "-X 'github.com/$(org)/$(target).Version=$(version)'" -o release/$(target)-darwin-amd64 $(gofile)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
