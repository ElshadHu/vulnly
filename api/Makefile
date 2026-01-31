.PHONY: build deploy local clean
build:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap ./cmd/api
deploy: build
	sam deploy
local:
	go run ./cmd/api
clean:
	rm -f bootstrap