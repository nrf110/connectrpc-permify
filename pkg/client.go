package connectpermify

import (
	"fmt"

	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
)

type CheckClient interface {
	Check(principal *Resource, attributes Attributes, config CheckConfig) (bool, error)
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

func (client *permifyCheckClient) Check(principal *Resource, attributes Attributes, config CheckConfig) (bool, error) {
	switch config.Type {
	case SINGLE:
		return client.check(principal, attributes, config)
	default:
		return false, fmt.Errorf("unexpected CheckType %s", config.Type)
	}
}

func (client *permifyCheckClient) check(principal *Resource, attributes Attributes, config CheckConfig) (bool, error) {
	request, err := config.Checks[0].toCheckRequest(principal, attributes)
	if err != nil {
		return false, err
	}
	return client.Client.Check(request)
}
