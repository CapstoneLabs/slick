# The default, used by Travis CI
test:
	./scripts/pre-commit.sh

build:
	env GO111MODULE=on go build ./...

get:
	env GO111MODULE=on go get ./...

cov: 
	env GO111MODULE=on go test -coverprofile=coverage.out 
	env GO111MODULE=on go tool cover -html=coverage.out

build-docs:
	cd docs-src && hugo

clean-docs:
	rm -rf docs/*

run-docs:
	cd docs-src && hugo server --watch