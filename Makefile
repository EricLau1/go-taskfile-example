build-proto:
	@protoc --go_out=. --go-grpc_out=. ./protos/*.proto
	@echo "proto done."

build-account:
	@go build -o bin/account account/main.go
	@echo "account done."

build-posts:
	@go build -o bin/posts posts/main.go
	@echo "posts done."

build-comments:
	@go build -o bin/comments comments/main.go
	@echo "comments done."

build: build-proto build-account build-posts build-comments
	@echo "build success."