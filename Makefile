.PHONY: build clean deploy

build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/alive src/alive.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

format:
	gofmt -w src/alive.go
