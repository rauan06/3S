# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
BINARY_NAME=triples
SOURCE_FILE=cmd/main.go

# Targets
all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(SOURCE_FILE)

run:
	$(GORUN) $(SOURCE_FILE)

start:
	./$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	clear
