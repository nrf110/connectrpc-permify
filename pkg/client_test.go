package connectpermit

import (
	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ovechkin-dm/mockio/mock"
)

var _ = Describe("Check", func() {
	const PERMISSION = "edit"

	var (
		stubPrincipal *Resource
		stubResource  *Resource
	)

	BeforeEach(func() {
		stubPrincipal = &Resource{
			ID: "1234",
		}
		stubResource = &Resource{
			Type: "widget",
			ID:   "abcde",
		}
	})

	When("config type is single", func() {
		It("should invoke Check", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			permitClient := mock.Mock[PermifyInterface]()
			mock.When(permitClient.Check(mock.Any[*permifypayload.PermissionCheckRequest]())).
				ThenReturn(true, nil)

			checkClient := NewCheckClient(permitClient)
			result, err := checkClient.Check(stubPrincipal, Attributes{}, CheckConfig{
				Type: SINGLE,
				Checks: []Check{
					{
						Permission: PERMISSION,
						Entity:     stubResource,
					},
				},
			})

			Expect(err).To(BeNil())
			Expect(result).To(BeTrue())

			mock.Verify(permitClient, mock.Once()).Check(mock.Any[*permifypayload.PermissionCheckRequest]())
		})
	})

	When("config type is public", func() {
		It("should return an error", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			permitClient := mock.Mock[PermifyInterface]()

			checkClient := NewCheckClient(permitClient)
			result, err := checkClient.Check(stubPrincipal, Attributes{}, CheckConfig{
				Type: PUBLIC,
				Checks: []Check{
					{
						Permission: PERMISSION,
						Entity:     stubResource,
					},
				},
			})

			Expect(err.Error()).To(Equal("unexpected CheckType public"))
			Expect(result).To(BeFalse())

			mock.Verify(permitClient, mock.Never()).Check(mock.Any[*permifypayload.PermissionCheckRequest]())
		})
	})
})
