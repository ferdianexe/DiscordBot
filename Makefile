.PHONY: all clean build run

# Default target
all: build

# Build the Go application
build:
	go build -o app

# Fetch dependencies and build the application
deps:
	go get -u github.com/mattn/go-sqlite3
	go build -o app

# Run the application
run:
	go build -o app
	./app

# Clean compiled binaries
clean:
	rm -f app
