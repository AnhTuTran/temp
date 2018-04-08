GOPATH = $(shell pwd)
GOSOURCE = $(shell ls **/*.go)
BINDIR = bin/
SOURCEDIR = src/
TARGET = mirage
GO = GOPATH=$(GOPATH) go
GIT = git

all: format $(BINDIR)$(TARGET)

format:
	./tools/go_source/recursive-gofmt.sh

golint:
	./tools/go_source/golint.sh

go-vet:
	./tools/go_source/go-vet.sh

$(BINDIR)$(TARGET):
	$(GO) build -ldflags="-s -w" -o $@ mirage

run: $(SOURCEDIR)$(TARGET).go
	$(GO) run $<

clean:
	rm -f $(BINDIR)$(TARGET)

git-push: format
	$(GIT) push

