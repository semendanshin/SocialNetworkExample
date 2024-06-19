# Alias for generating grpc code from proto files

# Generate grpc code from proto files
grpc:
	protoc -I api/proto api/proto/sso.proto --go_out=./gen/go/ --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative