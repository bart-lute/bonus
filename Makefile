build:
		go build -o ./target/bonus -ldflags "-s -w" ./cmd/bonus
		go build -o ./target/web -ldflags "-s -w" ./cmd/web

build-release:
		GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o bonus