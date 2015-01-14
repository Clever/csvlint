SHELL := /bin/bash
PKG = github.com/Clever/csvlint
SUBPKGS = $(addprefix $(PKG)/, cmd/csvlint)
PKGS = $(PKG) $(SUBPKGS)
VERSION := $(shell cat VERSION)
EXECUTABLE := csvlint
BUILDS := \
	build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64 \
	build/$(EXECUTABLE)-v$(VERSION)-linux-amd64 \
	build/$(EXECUTABLE)-v$(VERSION)-windows-amd64
COMPRESSED_BUILDS := $(BUILDS:%=%.tar.gz)
RELEASE_ARTIFACTS := $(COMPRESSED_BUILDS:build/%=release/%)

.PHONY: test golint

golint:
	@go get github.com/golang/lint/golint

test: $(PKGS)

$(PKGS): golint
	@go get -d -t $@
	@gofmt -w=true $(GOPATH)/src/$@/*.go
ifneq ($(NOLINT),1)
	@echo "LINTING..."
	@PATH=$(PATH):$(GOPATH)/bin golint $(GOPATH)/src/$@/*.go
	@echo ""
endif
ifeq ($(COVERAGE),1)
	@echo "TESTING WITH COVERAGE... $@"
	@go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	@go tool cover -html=$(GOPATH)/src/$@/c.out
else
	@echo "TESTING... $@"
	@go test $@ -test.v
endif

run:
	@go run cmd/csvlint/main.go

build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build -o "$@/$(EXECUTABLE)" $(PKG)/cmd/csvlint
build/$(EXECUTABLE)-v$(VERSION)-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -o "$@/$(EXECUTABLE)" $(PKG)/cmd/csvlint
build/$(EXECUTABLE)-v$(VERSION)-windows-amd64:
	GOARCH=amd64 GOOS=windows go build -o "$@/$(EXECUTABLE).exe" $(PKG)/cmd/csvlint
build: $(BUILDS)
%.tar.gz: %
	tar -C `dirname $<` -zcvf "$<.tar.gz" `basename $<`
$(RELEASE_ARTIFACTS): release/% : build/%
	mkdir -p release
	cp $< $@
release: $(RELEASE_ARTIFACTS)

clean:
	rm -rf build release
