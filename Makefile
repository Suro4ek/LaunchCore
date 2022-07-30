proto server:
	protoc --go_out=. \
                   --go-grpc_out=.  \
                   proto/server.proto   
proto user:
	protoc --go_out=. \
                   --go-grpc_out=. \
                   proto/user.proto   