package connectpermify

import (
	"context"

	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	"google.golang.org/protobuf/types/known/structpb"
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
	Depth      *int32
}

func (c Check) toCheckRequest(
	ctx context.Context,
	principal *Resource,
	attributes Attributes,
	schemaVersion string,
) (*permifypayload.PermissionCheckRequest, error) {
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
		Metadata: &permifypayload.PermissionCheckRequestMetadata{
			SchemaVersion: schemaVersion,
			SnapToken:     GetSnapToken(ctx),
			Depth:         *c.Depth,
		},
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
	IsPublic bool
	Checks   []Check
}

type Checkable interface {
	GetChecks() CheckConfig
}
