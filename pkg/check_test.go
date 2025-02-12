package connectpermit

import (
	permifypayload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ = Describe("toCheckRequest", func() {
	When("attributes are present", func() {
		It("should include them on the CheckRequest", func(ctx SpecContext) {
			attributes := Attributes{
				"foo": "bar",
			}
			user := Resource{
				ID:   "abcde",
				Type: "user",
			}

			check := Check{
				Permission: "edit",
				Entity: &Resource{
					Type: "Widget",
					ID:   "1234",
				},
			}
			req, err := check.toCheckRequest(&user, attributes)
			Expect(err).To(BeNil())

			data, err := structpb.NewStruct(attributes)
			Expect(req).To(BeEquivalentTo(&permifypayload.PermissionCheckRequest{
				TenantId: "t1",
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
			}))
		})
	})
})
