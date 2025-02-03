package connectpermit

import (
	"github.com/auth0/go-jwt-middleware/v2/validator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultClaimsMapper", func() {
	var customClaims testCustomClaims
	var claims *validator.ValidatedClaims

	BeforeEach(func() {
		customClaims = testCustomClaims{
			Roles:          []string{"admin"},
			OrganizationID: "acme",
		}
		claims = &validator.ValidatedClaims{
			RegisteredClaims: validator.RegisteredClaims{},
			CustomClaims:     &customClaims,
		}
	})

	It("should set the subject as the key", func(ctx SpecContext) {
		mapper := DefaultClaimsMapper(DefaultCustomClaimsMapper[*testCustomClaims]())
		user, err := mapper(claims)
		Expect(err).To(BeNil())
		Expect(user.Key).To(Equal(claims.RegisteredClaims.Subject))
	})

	It("should convert the custom claims to a map of attributes", func(ctx SpecContext) {
		mapper := DefaultClaimsMapper(DefaultCustomClaimsMapper[*testCustomClaims]())
		user, err := mapper(claims)
		Expect(err).To(BeNil())
		Expect(user.Attributes).To(HaveKeyWithValue("Roles", []string{"admin"}))
		Expect(user.Attributes).To(HaveKeyWithValue("OrganizationID", "acme"))
	})
})
