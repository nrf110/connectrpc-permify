package connectpermify

import (
	"context"

	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
)

type CheckClient interface {
	Check(ctx context.Context, principal *Resource, attributes Attributes, config CheckConfig) (bool, error)
}

type PermifyInterface interface {
	Check(request *permifypayload.PermissionCheckRequest) (bool, error)
}

type permifyCheckClient struct {
	Client PermifyInterface
}

func NewCheckClient(client PermifyInterface) CheckClient {
	return &permifyCheckClient{
		Client: client,
	}
}

func (client *permifyCheckClient) Check(ctx context.Context, principal *Resource, attributes Attributes, config CheckConfig) (bool, error) {
	// TODO: parallize this until the permify client supports a native bulk operation
	results := make([]bool, len(config.Checks))
	for idx, check := range config.Checks {
		result, err := client.check(principal, attributes, check)
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

func (client *permifyCheckClient) check(principal *Resource, attributes Attributes, check Check) (bool, error) {
	request, err := check.toCheckRequest(principal, attributes)
	if err != nil {
		return false, err
	}
	return client.Client.Check(request)
}
