protoc --go_out=. *.proto
protoc --go-grpc_out=. *.proto

mv protoc_grpc.pb.go protoc_grpc.go
mv protoc.pb.go protoc.go