.PHONY: build clean

build: clean
	go test
	go build

clean:
	rm -f "$(basename $(pwd))"
