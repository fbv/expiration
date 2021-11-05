projects=\
	example

all: build

clean:
	@rm -rf bin/*

build: $(projects)

$(projects):
	go build -o bin/$@ cmd/$@/main.go
	@#GOOS=darwin GOARCH=arm64 go build -o bin/$@-macos-arm64 cmd/$@/main.go

test:
	go test ./...

.PHONY: all clean build test $(projects)
