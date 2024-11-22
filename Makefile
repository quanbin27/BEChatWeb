run:
	@go run cmd/api_gateway/main.go
	@go run cmd/grpc_server/main.go

gen_user:
	@protoc \
		--proto_path=protobuf "protobuf/users.proto" \
		--go_out=services/common/genproto/users --go_opt=paths=source_relative \
  	--go-grpc_out=services/common/genproto/users --go-grpc_opt=paths=source_relative
gen_group:
	@protoc \
    		--proto_path=protobuf "protobuf/groups.proto" \
    		--go_out=services/common/genproto/groups --go_opt=paths=source_relative \
      	--go-grpc_out=services/common/genproto/groups --go-grpc_opt=paths=source_relative
gen_message:
	@protoc \
    		--proto_path=protobuf "protobuf/messages.proto" \
    		--go_out=services/common/genproto/messages --go_opt=paths=source_relative \
      	--go-grpc_out=services/common/genproto/messages --go-grpc_opt=paths=source_relative