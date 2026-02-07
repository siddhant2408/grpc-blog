# Architecture Overview

This project follows clean architecture principles:

- cmd/: application entrypoints (server, client)
- internal/app/: business logic
- internal/transport/: gRPC adapters
- internal/infra/: logging, tracing
- internal/container/: dependency injection (Uber Dig)

Observability:
- Structured logging with Zap
- Distributed tracing with OpenTelemetry

All dependencies are injected and no global state is used.
