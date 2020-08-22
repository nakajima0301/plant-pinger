GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=ppinger
BINARY_WINDOWS=$(BINARY_NAME).exe

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WINDOWS)

# Cross-compile
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v