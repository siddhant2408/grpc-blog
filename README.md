# gRPC Blog Service

This repository contains a production-grade gRPC service written in Go.

## Features
- CRUD operations for blog posts
- gRPC API
- In-memory storage
- Structured logging (Zap)
- Distributed tracing (OpenTelemetry)
- Dependency Injection (Uber Dig)
- High unit test coverage

## Repository Structure
cmd/                Application entrypoints
internal/app/       Business logic
internal/transport/ gRPC adapters
internal/infra/     Logging & tracing
internal/container/ Dependency injection
proto/              Protobuf definitions
docs/               Documentation

## Setup

### Prerequisites
- Go 1.21+
- protoc
- protoc-gen-go
- protoc-gen-go-grpc

### Install dependencies
go mod tidy

### Generate protobuf code
protoc   --go_out=.   --go-grpc_out=.   --go_opt=paths=source_relative   --go-grpc_opt=paths=source_relative   proto/blog.proto

## Run server
go run cmd/server/main.go

## Run client
go run cmd/client/main.go

## Run tests
go test ./...

## Documentation
- docs/api.md
- docs/architecture.md
