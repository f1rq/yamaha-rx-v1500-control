BINARY=yamaha-rx-v1500-control

build-rpi:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY) .

clean:
	rm -f $(BINARY)