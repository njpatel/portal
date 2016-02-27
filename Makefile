VERSION = 0.0.1

all: portal

portal:
	@go build -ldflags "-X main.version=${VERSION}"

# gen generates the protobuf source files from the .proto definition files
gen:
	@protoc --go_out=plugins=grpc:. ./api/portal.proto

test: portal
	@go test -v

.PHONY: portal gen test
