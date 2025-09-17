package connectpermify

import (
	"testing"

	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestToCheckRequest(t *testing.T) {
	depth := int32(10)

	t.Run("merge check attributes with claim attributes", func(t *testing.T) {
		attributes := Attributes{
			"foo": "bar",
		}
		user := Resource{
			ID:   "abcde",
			Type: "user",
		}

		check := Check{
			TenantID:   "t2",
			Permission: "edit",
			Entity: &Resource{
				Type:       "Widget",
				ID:         "1234",
				Attributes: Attributes{"quux": "corge"},
			},
			Depth: &depth,
		}
		req, err := check.toCheckRequest(t.Context(), &user, attributes, "")
		assert.NoError(t, err)

		data, err := structpb.NewStruct(Attributes{
			"foo":  "bar",
			"quux": "corge",
		})
		assert.NoError(t, err)

		expected := &permifypayload.PermissionCheckRequest{
			TenantId: "t2",
			Subject: &permifypayload.Subject{
				Id:   user.ID,
				Type: user.Type,
			},
			Permission: "edit",
			Entity: &permifypayload.Entity{
				Type: "Widget",
				Id:   "1234",
			},
			Context: &permifypayload.Context{
				Data: data,
			},
			Metadata: &permifypayload.PermissionCheckRequestMetadata{
				Depth: depth,
			},
		}
		assert.Equal(t, expected, req)
	})
}
