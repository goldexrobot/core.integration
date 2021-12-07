//go:build tools
// +build tools

package main

import (
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "golang.org/x/tools/cmd/stringer"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
