gen:
	@protoc \
		--proto_path=protobuf "protobuf/users.proto" \
		--go_out=services/common/genproto/users --go_opt=paths=source_relative \
  	--go-grpc_out=services/common/genproto/users --go-grpc_opt=paths=source_relative