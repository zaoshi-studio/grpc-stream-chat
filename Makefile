
projectDir = $(shell pwd)

.PHONY: pb
pb:
	protoc  --proto_path=$(projectDir)/proto --go_out=$(projectDir)/pb --go_opt=paths=source_relative \
			--go-grpc_out=$(projectDir)/pb --go-grpc_opt=paths=source_relative \
			$(projectDir)/proto/protocol/*.proto