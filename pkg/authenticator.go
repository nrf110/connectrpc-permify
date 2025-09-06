package connectpermify

import (
	"context"

	"connectrpc.com/connect"
)

type AuthenticationResult struct {
	Context    context.Context
	Principal  *Resource
	Attributes Attributes
}

type Authenticator interface {
	Authenticate(ctx context.Context, req connect.AnyRequest) (*AuthenticationResult, error)
}
