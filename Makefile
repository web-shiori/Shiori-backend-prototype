.PHONY: build clean deploy

build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/alive src/alive.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/post src/post.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/list src/list.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

delete:
	sls remove

format:
	gofmt -w src/alive.go
	gofmt -w src/post.go
