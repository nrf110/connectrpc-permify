package connectpermify

import (
	"context"

	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	"google.golang.org/grpc"
)

type CheckClient interface {
	Check(ctx context.Context, principal *Resource, attributes Attributes, config CheckConfig) (bool, error)
}

type PermifyInterface interface {
	Check(ctx context.Context, request *permifypayload.PermissionCheckRequest, opts ...grpc.CallOption) (*permifypayload.PermissionCheckResponse, error)
}

type permifyCheckClient struct {
	schemaVersion string
	defaultDepth  int32
	client        PermifyInterface
}

type ClientOption func(*permifyCheckClient)

func WithDefaultDepth(depth int32) ClientOption {
	return func(client *permifyCheckClient) {
		client.defaultDepth = depth
	}
}

func WithSchemaVersion(version string) ClientOption {
	return func(client *permifyCheckClient) {
		client.schemaVersion = version
	}
}

func NewCheckClient(client PermifyInterface, options ...ClientOption) CheckClient {
	c := &permifyCheckClient{
		client:       client,
		defaultDepth: 10,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (client *permifyCheckClient) Check(ctx context.Context, principal *Resource, attributes Attributes, config CheckConfig) (bool, error) {
	// TODO: parallize this until the permify client supports a native bulk operation
	results := make([]bool, len(config.Checks))
	for idx, check := range config.Checks {
		result, err := client.check(ctx, principal, attributes, check)
		if err != nil {
			return false, err
		}
		results[idx] = result
	}

	result := true
	for _, r := range results {
		result = result && r
	}
	return result, nil
}

func (client *permifyCheckClient) check(
	ctx context.Context,
	principal *Resource,
	attributes Attributes,
	check Check,
) (bool, error) {
	if check.Depth == nil {
		check.Depth = &client.defaultDepth
	}
	request, err := check.toCheckRequest(ctx, principal, attributes, client.schemaVersion)
	if err != nil {
		return false, err
	}
	res, err := client.client.Check(ctx, request)
	if err != nil {
		return false, err
	}

	return res.Can == permifypayload.CheckResult_CHECK_RESULT_ALLOWED, nil
}
