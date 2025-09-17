# ConnectRPC Permify Interceptor

A Go library that provides authorization middleware for ConnectRPC services using [Permify](https://permify.co) for fine-grained access control.

## Overview

This interceptor validates JWT tokens and performs permission checks against Permify before allowing requests to proceed. It extracts authorization information from JWT claims and maps them to Permify's resource-based permission model.

## Features

- **Pluggable Authentication**: Implement the Authenticator interface to provide a security principal to authorize against
- **Protocol Buffer Extensions**: Annotate your protobuf definitions to automatically configure permissions
- **Flexible Authorization**: Support for public endpoints and disabled authorization modes
- **ConnectRPC Integration**: Seamlessly integrates as a ConnectRPC interceptor

## Installation

```bash
go get github.com/nrf110/connectrpc-permify
```

## Quick Start

### 1. Annotate Your Protobuf Services

Use the custom extensions to mark your protobuf messages and methods:

```protobuf
import "nrf110/permify/v1/permify.proto";

message GetUserRequest {
  string user_id = 1 [(nrf110.permify.v1.resource_id) = true];
}

message User {
  option (nrf110.permify.v1.resource_type) = "user";
  string id = 1;
  string email = 2;
}

service UserService {
  rpc GetUser(GetUserRequest) returns (User) {
    option (nrf110.permify.v1.action) = "read";
    option (nrf110.permify.v1.depth) = 2;  // Optional
  }

  rpc PublicHealthCheck(HealthCheckRequest) returns (HealthCheckResponse) {
    option (nrf110.permify.v1.public) = true;  // No auth required
  }
}
```

### 2. Set Up the Interceptor

```go
package main

import (
    "github.com/nrf110/connectrpc-permify"
    permifyclient "buf.build/gen/go/permifyco/permify/connectrpc/go/grpc/v1/grpcv1connect"
)

func main() {
    // Configure Permify client
    permify := permifyclient.NewPermissionClient(
        http.DefaultClient,
        "http://localhost:3476", // Permify server URL
    )

    // An Authenticator defines the logic for finding credentials in the request, verifying them, and returning the Principal
    // NewCustomAuthenticator isn't a real method.  Separate libraries defining Authenticator implementations are available.
    authenticator := NewCustomAuthenticator()

    // Create the interceptor
    interceptor := connectpermify.NewPermifyInterceptor(
        connectpermify.NewCheckClient(permify),
        authenticator,
        func() bool { return true }, // Enable/disable authorization.  Useful for testing locally without checks, or while migrating from another authorization solution.
    )

    // Use with your ConnectRPC server
    mux := http.NewServeMux()
    path, handler := yourv1connect.NewUserServiceHandler(
        &userService{},
        connect.WithInterceptors(interceptor),
    )
    mux.Handle(path, handler)
}
```

### 3. Implement Required Interfaces

Your request messages must implement the `Checkable` interface:

```go
// Usually auto-generated from protobuf annotations
func (r *GetUserRequest) GetChecks() connectpermify.CheckConfig {
    return connectpermify.CheckConfig{
        Checks: []connectpermify.Check{
            {
                TenantID:   r.GetTenantId(), // If using multi-tenancy
                Permission: "read",
                Entity: &connectpermify.Resource{
                    Type: "user",
                    ID:   r.GetUserId(),
                },
                Depth: int32(10),
            },
        },
    }
}
```

## Protocol Buffer Extensions

The library provides custom protobuf extensions for automatic permission configuration:

- `(nrf110.permify.v1.resource_id) = true` - Mark fields that contain resource IDs
- `(nrf110.permify.v1.tenant_id) = true` - Mark fields that contain tenant IDs
- `(nrf110.permify.v1.attribute_name) = "attr"` - Map fields to Permify attributes
- `(nrf110.permify.v1.resource_type) = "user"` - Set resource type for messages
- `(nrf110.permify.v1.action) = "read"` - Set required permission for methods
- `(nrf110.permify.v1.public) = true` - Mark methods as public (no auth required)
- `(nrf110.permify.v1.depth) = true` - Set the depth parameter for the check request, in the event the graph of permissions gets into a loop

## Check Client Configuration

```go
checkClient := connectpermify.NewCheckClient(
  permifyClient,

  // Optional.  Defaults to 10
  connectpermify.WithDefaultDepth(10),

  // Optional.  Defaults to ""
  connectpermify.WithSchemaVersion("latest"),
)
```

## Development

### Prerequisites

- Go 1.24+
- [Buf CLI](https://buf.build/docs/installation)
- [Permify server](https://docs.permify.co/getting-started/installation) for testing

### Building

```bash
# Generate protobuf code
make gen

# Run tests
make test

# Clean generated files
make clean

# Update dependencies
make update
```

### Development Workflow

1. **Modify protobuf definitions** in `proto/nrf110/permify/v1/`
2. **Regenerate code** with `make gen`
3. **Run tests** with `make test`
4. **Validate changes** against existing services

### Testing

The library uses Go's standard testing framework with [testify](https://github.com/stretchr/testify) for assertions and [mockio](https://github.com/ovechkin-dm/mockio) for mocking:

```bash
# Run all tests
make test

# Run tests for specific package
go test -v ./pkg/
```

## License

This project is licensed under MIT license.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run `make test` to validate
5. Submit a pull request
