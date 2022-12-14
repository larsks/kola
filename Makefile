TARGET = kola
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
TEMPLATES = $(shell find . -type f -name '*.tpl' -not -path "./vendor/*")

VERSION = $(shell git describe --tags --exact-match 2> /dev/null || echo development)
COMMIT = $(shell git rev-parse --short=10 HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")

DESTDIR=
prefix=$(HOME)
bindir=$(prefix)/bin
INSTALL=install

BUILDDATA = \
	-X "kola/version.Version=$(VERSION)" \
	-X "kola/version.BuildDate=$(DATE)" \
	-X "kola/version.BuildRef=$(COMMIT)"

LDFLAGS = -ldflags '$(BUILDDATA)'

all: $(TARGET)

$(TARGET): .checked $(SRC) $(TEMPLATES) go.sum
	go build $(LDFLAGS) -o $(TARGET)

foo:
	echo $(TEMPLATES)

check: .checked

.checked: $(SRC) go.sum
	golangci-lint run | tee $@

go.sum: go.mod
	go mod tidy && touch $@

install: $(TARGET)
	$(INSTALL) -m 755 $(TARGET) $(DESTDIR)$(bindir)/$(TARGET)

clean:
	rm -f $(TARGET) .checked
