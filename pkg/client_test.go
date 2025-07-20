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

	t.Run("invokes Check when config type is single", func(t *testing.T) {
		mock.SetUp(t)
		permifyClient := mock.Mock[PermifyInterface]()
		mock.When(permifyClient.Check(mock.Any[*permifypayload.PermissionCheckRequest]())).
			ThenReturn(true, nil)

		checkClient := NewCheckClient(permifyClient)
		result, err := checkClient.Check(stubPrincipal, Attributes{}, CheckConfig{
			Type: SINGLE,
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

	t.Run("returns an error when config type is public", func(t *testing.T) {
		mock.SetUp(t)
		permifyClient := mock.Mock[PermifyInterface]()

		checkClient := NewCheckClient(permifyClient)
		result, err := checkClient.Check(stubPrincipal, Attributes{}, CheckConfig{
			Type: PUBLIC,
			Checks: []Check{
				{
					Permission: PERMISSION,
					Entity:     stubResource,
				},
			},
		})

		assert.EqualError(t, err, "unexpected CheckType public")
		assert.False(t, result)

		mock.Verify(permifyClient, mock.Never()).Check(mock.Any[*permifypayload.PermissionCheckRequest]())
	})
}
