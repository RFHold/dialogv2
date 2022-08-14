//go:build tools
// +build tools

package build

import (
	_ "github.com/pressly/goose"
	_ "golang.org/x/lint/golint"
	_ "google.golang.org/grpc"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
