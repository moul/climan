GOPKG ?=	moul.io/climan

include rules.mk

generate: install
	GO111MODULE=off go get github.com/campoy/embedmd
	mkdir -p .tmp
	go doc -all > .tmp/godoc.txt
	embedmd -w README.md
	rm -rf .tmp
.PHONY: generate
