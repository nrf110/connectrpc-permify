package connectpermify

import (
	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

type CheckType string

const (
	SINGLE CheckType = "single"
	PUBLIC CheckType = "public"
)

type Attributes map[string]any

type Resource struct {
	Type       string
	ID         string
	Attributes Attributes
}

type Check struct {
	TenantID   string
	Permission string
	Entity     *Resource
}

func (c Check) toCheckRequest(principal *Resource, attributes Attributes) (*permifypayload.PermissionCheckRequest, error) {
	tenantId := "t1"
	if c.TenantID != "" {
		tenantId = c.TenantID
	}

	subject := &permifypayload.Subject{
		Type: principal.Type,
		Id:   principal.ID,
	}

	entity := &permifypayload.Entity{
		Type: c.Entity.Type,
		Id:   c.Entity.ID,
	}

	req := &permifypayload.PermissionCheckRequest{
		TenantId:   tenantId,
		Subject:    subject,
		Permission: c.Permission,
		Entity:     entity,
	}

	mergedAttributes := make(Attributes)
	for k, v := range attributes {
		mergedAttributes[k] = v
	}
	for k, v := range c.Entity.Attributes {
		mergedAttributes[k] = v
	}
	if len(mergedAttributes) > 0 {
		data, err := structpb.NewStruct(mergedAttributes)
		if err != nil {
			return nil, err
		}
		req.Context = &permifypayload.Context{
			Data: data,
		}
	}

	return req, nil
}

type CheckConfig struct {
	Type   CheckType
	Checks []Check
}

func (config CheckConfig) IsPublic() bool {
	return config.Type == PUBLIC
}

type Checkable interface {
	GetChecks() CheckConfig
}
