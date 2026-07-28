.PHONY: build test vet fmt fmt-check check clean

build:
	go build ./...

test:
	go test ./... -race -cover

vet:
	go vet ./...

fmt:
	gofmt -l -w .

fmt-check:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)

check: fmt-check vet test

clean:
	rm -rf bin dist
