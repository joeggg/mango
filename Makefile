.PHONY: test

test:
	go test -coverprofile cover.out -v ./...
	go tool cover --func cover.out

proto:
	protoc -I=pb/raw --go_out=pb pb/raw/*
