PROTO_SRC := api/proto
PROTO_OUT := api/gen/go

gen:
	protoc \
		--go_out=api/gen/go \
		--go-grpc_out=api/gen/go \
		--go_opt=paths=import \
		--go-grpc_opt=paths=import \
		api/proto/*.proto
