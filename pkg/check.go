package connectpermify

import (
	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

type CheckType string

const (
	SINGLE CheckType = "single"
	PUBLIC           = "public"
)

type Attributes map[string]any

type Resource struct {
	Type string
	ID   string
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

	if attributes != nil {
		data, err := structpb.NewStruct(attributes)
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
