TARGET = kola
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

VERSION = $(shell git describe --tags --exact-match 2> /dev/null || echo development)
COMMIT = $(shell git rev-parse --short=10 HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")

BUILDDATA = \
	-X "kola/version.Version=$(VERSION)" \
	-X "kola/version.BuildDate=$(DATE)" \
	-X "kola/version.BuildRef=$(COMMIT)"

LDFLAGS = -ldflags '$(BUILDDATA)'

all: $(TARGET)

$(TARGET): $(SRC)
	go build $(LDFLAGS) -o $(TARGET)

clean:
	rm -f $(TARGET)
