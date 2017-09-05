.PHONY: build clean install

build: clean
	go test
	cd cmd && go test
	cd v1 && go test
	cd v1/cmd && go test
	go build

clean:
	rm -f "$(basename $(pwd))"

install: build
	go install github.com/ynishi/publish/cmd/publish
