.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go build -o ./.bin/pcstore ./cmd/pcstore/

run: build
	./.bin/pcstore