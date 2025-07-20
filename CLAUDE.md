# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a ConnectRPC interceptor library for Permify authorization written in Go. It provides middleware that validates JWT tokens and performs permission checks against Permify before allowing requests to proceed.

## Common Development Commands

### Building and Code Generation
```bash
make gen          # Clean, update deps, and generate protobuf code
make clean        # Remove generated files (gen/ and dist/)
make update       # Run go mod tidy
```

### Testing
```bash
make test         # Run all tests (runs make gen first)
go test ./pkg/    # Run tests directly without regeneration
```

### Protocol Buffer Development
```bash
buf generate      # Generate Go code from protobuf definitions
buf lint          # Lint protobuf files
buf breaking      # Check for breaking changes
```

## Code Architecture

### Core Components

**PermifyInterceptor** (`pkg/interceptor.go`) - Main authorization middleware that orchestrates:
- Token extraction from Authorization headers
- JWT validation using OIDC
- Claims mapping to Permify resources/attributes  
- Permission checks against Permify API

**Token Management**:
- `TokenExtractor` - Extracts Bearer tokens
- `TokenValidator` - OIDC JWT validation with configurable issuers/audiences

**Authorization Flow**:
- `ClaimsMapper` - Maps JWT claims to Permify entities (supports generics)
- `CheckClient` - Abstraction for Permify permission API calls
- `CheckConfig` - Permission requirements configuration

**Protocol Extensions** (`proto/nrf110/permify/v1/`):
- Custom protobuf extensions for marking resource IDs, tenant IDs, and actions
- Used for automatic permission inference from service definitions

### Key Patterns

- **Interface-based design** - All core components use interfaces for testability
- **Dependency injection** - Components are injected rather than instantiated
- **Middleware pattern** - Follows ConnectRPC interceptor conventions
- **Configuration-driven** - Supports public endpoints and disabled auth modes

### Important Files

- `pkg/interceptor.go` - Main authorization middleware implementation
- `pkg/client.go` - Permify client abstraction and configuration
- `pkg/token-validator.go` - OIDC JWT validation logic  
- `pkg/claims-mapper.go` - JWT claims to Permify mapping
- `proto/nrf110/permify/v1/permify.proto` - Custom protobuf extensions

### Testing Strategy

- Uses Go stdlib testing with testify assertions
- Mockio for generating mocks
- Each component has corresponding `_test.go` file
- Tests validate authorization flows, token validation, and claims mapping

## Development Environment

The project includes devcontainer support with configurations for VS Code and GoLand. The container includes Go, Buf, and all necessary tools for development.