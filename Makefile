generate_grpc_code:
	protoc \
    --go_out=converter \
    --go_opt=paths=source_relative \
    --go-grpc_out=converter \
    --go-grpc_opt=paths=source_relative \
    converter.proto
