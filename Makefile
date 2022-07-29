proto server:
	protoc --go_out=. --go_opt=paths=source_relative \
                   --go-grpc_out=.  \
                   proto/server.proto   
proto user:
	protoc --go_out=. --go_opt=paths=source_relative \
                   --go-grpc_out=. \
                   proto/user.proto   