mkdir -p ./proto
protoc --proto_path=/proto/ --go-grpc_out=proto --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative --go_out=proto fibonacci.proto
