VERSION=0.1.0

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
BINARY=ghworm
BINARY_LNX=$(BINARY)
BINARY_OSX=$(BINARY)
BINARY_WIN=$(BINARY).exe
RELEASEDIR=release
ORG=brightm
REPO=$(BINARY)

default: test build fmt

build:
	$(GOFMT) *.go
	$(GOBUILD) -o $(BINARY) -v

fmt:
	$(GOFMT) *.go
	$(GOFMT) cmd/*.go

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(RELEASEDIR)
	git clean -fxd

run:
	$(GOBUILD) -o $(BINARY) -v ./...

deps:
	$(GOGET) -u github.com/google/go-github/github
	$(GOGET) -u golang.org/x/oauth2
	$(GOGET) -u github.com/tcnksm/ghr
	$(GOGET) -u github.com/spf13/cobra/cobra
	$(GOGET) -u github.com/ianmcmahon/encoding_ssh

cross:
	mkdir -p $(RELEASEDIR)/lnx $(RELEASEDIR)/osx $(RELEASEDIR)/win
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(RELEASEDIR)/lnx/$(BINARY_LNX) -v
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(RELEASEDIR)/osx/$(BINARY_OSX) -v
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(RELEASEDIR)/win/$(BINARY_WIN) -v

packaging: cross
	cp README.md $(RELEASEDIR)/lnx
	tar -zcf $(RELEASEDIR)/$(BINARY)-lnx.tar.gz -C $(RELEASEDIR)/lnx .
	cp README.md $(RELEASEDIR)/osx
	tar -zcf $(RELEASEDIR)/$(BINARY)-osx.tar.gz -C $(RELEASEDIR)/osx .
	cp README.md $(RELEASEDIR)/win
	zip -r $(RELEASEDIR)/$(BINARY)-win.zip -j $(RELEASEDIR)/win

release: packaging
	ghr \
	-u $(ORG) \
	-r $(REPO) \
	-b "send secrets over gists with GitHub ssh keys" \
	$(VERSION) $(RELEASEDIR)
