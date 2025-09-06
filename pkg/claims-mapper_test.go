package connectpermify

import (
	"testing"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/stretchr/testify/assert"
)

func TestDefaultClaimsMapper(t *testing.T) {
	var (
		customClaims = testCustomClaims{
			Roles:          []string{"admin"},
			OrganizationID: "acme",
		}
		claims = &validator.ValidatedClaims{
			RegisteredClaims: validator.RegisteredClaims{},
			CustomClaims:     &customClaims,
		}
	)

	t.Run("should set the subject as the key", func(t *testing.T) {
		mapper := DefaultClaimsMapper(DefaultCustomClaimsMapper[*testCustomClaims]())
		principal, _, err := mapper(claims)
		assert.NoError(t, err)
		assert.Equal(t, claims.RegisteredClaims.Subject, principal.ID)
	})

	t.Run("should convert the custom claims to a map of attributes", func(t *testing.T) {
		mapper := DefaultClaimsMapper(DefaultCustomClaimsMapper[*testCustomClaims]())
		principal, attributes, err := mapper(claims)
		assert.NoError(t, err)
		assert.Equal(t, claims.RegisteredClaims.Subject, principal.ID)
		assert.Equal(t, attributes["Roles"], []string{"admin"})
		assert.Equal(t, attributes["OrganizationID"], "acme")
	})
}
