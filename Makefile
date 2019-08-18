build:
	go build -v

test:
	go test -p 1 ./... -cover

infratest:
	INFRA_TESTS=yes go test -p 1 ./... -cover

buildexample:
	docker build -t service -f Dockerfile.service .
