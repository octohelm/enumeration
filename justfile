tidy: gen fmt
    go mod tidy

test:
    CGO_ENABLED=0 go test -failfast -v ./...

test-race:
    CGO_ENABLED=1 go test -v -race ./...

fmt:
    go tool fmt .

dep:
    go get -u ./...

gen:
    go tool gen
