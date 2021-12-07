PROTOC_INC := -I ./../proto -I /usr/local/include/ -I ${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@$(shell go list -m -f '{{ .Version }}' "github.com/envoyproxy/protoc-gen-validate")

all: file.proto
	@:

%.proto:
	@mkdir -p $*
	protoc $(PROTOC_INC) --go_out=./$*/ --go_opt=paths=source_relative --go-grpc_out=paths=source_relative:./$*/ --validate_out=lang=go,paths=source_relative:./$*/ file.proto
