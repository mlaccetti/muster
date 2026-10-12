default:
    @just --list

fmt:
    gofumpt -w .

lint:
    golangci-lint run ./...

test:
    go test ./... -race -count=1

check-no-sleep:
    ./scripts/check-no-sleep.sh

build:
    go build -o dist/muster ./cmd/muster

# Run the CLI without installing (pass flags after the recipe name).
run *args:
    go run ./cmd/muster {{args}}

install:
    go install ./cmd/muster

check: fmt lint check-no-sleep test

release-dry:
    goreleaser release --snapshot --clean
