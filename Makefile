# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=gcs
PID:=$(shell cat ./gcs.pid)
LOG_FILE=./log/gcs.log

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_NAME)
	#rm -rf $(LOG_FILE)
	kill -9 $(PID)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
update:
	python scripts/updateset.py
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME) -mode update