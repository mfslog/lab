.PHONY: all clean

OUTPUT=receivetime

all: clean
	go build -o ./bin/${OUTPUT} cmd/main.go

clean:
	rm -f bin/${OUTPUT}

docker:
	docker build -t receivetime .