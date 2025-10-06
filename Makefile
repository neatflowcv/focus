.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o . ./cmd/...

.PHONY: install
install:
	CGO_ENABLED=0 GOOS=linux go install ./cmd/...

.PHONY: docs
docs:
	goa gen github.com/neatflowcv/focus/design

.PHONY: update
update:
	go get -u -t ./...
	go mod tidy
	go mod vendor

.PHONY: fix
fix:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0
	golangci-lint run --fix

.PHONY: lint
lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0
	golangci-lint run --allow-parallel-runners

.PHONY: test
test:
	go test -race -shuffle=on ./...

.PHONY: cover
cover:
	go test ./... --coverpkg ./... -coverprofile=c.out
	go tool cover -html="c.out"
	rm c.out