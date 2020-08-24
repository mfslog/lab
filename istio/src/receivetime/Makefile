.PHONY: all clean

OUTPUT=receivetime

all: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${OUTPUT} cmd/main.go

clean:
	rm -f bin/${OUTPUT}

docker:
	docker build --build-arg git_commit=$(git rev-parse --short HEAD) -t receivetime:latest .