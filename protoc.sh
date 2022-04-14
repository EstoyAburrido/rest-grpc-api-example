mkdir -p ./app/grpc/proto
protoc --proto_path=proto/ --go-grpc_out=app/grpc/proto --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative --go_out=app/grpc/proto fibonacci.proto
