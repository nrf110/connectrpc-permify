package connectpermify

import (
	"testing"

	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	const PERMISSION = "edit"

	var (
		stubPrincipal = &Resource{
			ID: "1234",
		}
		stubResource = &Resource{
			Type: "widget",
			ID:   "abcde",
		}
	)

	t.Run("invokes Check when not public", func(t *testing.T) {
		mock.SetUp(t)
		permifyClient := mock.Mock[PermifyInterface]()
		mock.When(permifyClient.Check(mock.Any[*permifypayload.PermissionCheckRequest]())).
			ThenReturn(true, nil)

		checkClient := NewCheckClient(permifyClient)
		result, err := checkClient.Check(t.Context(), stubPrincipal, Attributes{}, CheckConfig{
			Checks: []Check{
				{
					Permission: PERMISSION,
					Entity:     stubResource,
				},
			},
		})

		assert.NoError(t, err)
		assert.True(t, result)
		mock.Verify(permifyClient, mock.Once()).Check(mock.Any[*permifypayload.PermissionCheckRequest]())
	})
}
