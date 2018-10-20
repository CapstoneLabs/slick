# The default, used by Travis CI
test:
	./scripts/pre-commit.sh

build:
	go build ./...

get:
	env GO111MODULE=on go get ./...

cov: 
	go test -coverprofile=coverage.out 
	go tool cover -html=coverage.out