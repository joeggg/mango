.PHONY: test

test:
	go test -coverprofile cover.out -v

proto:
	protoc -I=pb/raw --go_out=pb pb/raw/*
