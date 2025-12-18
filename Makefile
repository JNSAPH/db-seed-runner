EXECUTABLE = db-seed-runner
BUILD_DIR = bin
RELEASE_DIR = release

build:
	go build -o bin/$(EXECUTABLE) 

run: build
	./bin/$(EXECUTABLE)