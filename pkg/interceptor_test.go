package connectpermit

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ovechkin-dm/mockio/mock"
)

var _ = Describe("Authorizing a ConnectRPC request", func() {
	var claims *validator.ValidatedClaims

	BeforeEach(func() {
		claims = &validator.ValidatedClaims{
			RegisteredClaims: validator.RegisteredClaims{
				Subject: "abcde",
			},
		}
	})

	When("the request is Checkable", func() {
		It("should invoke the next handler when the check call returns true", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(true, nil)

			tokenValidator := mock.Mock[TokenValidator]()
			mock.When(tokenValidator.Validate(mock.AnyContext(), mock.AnyString())).ThenReturn(claims, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, alwaysEnabled)
			result, err := interceptor(next)(ctx, req)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(res))
		})

		It("should invoke the next handler when the CheckConfig is public", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(true, nil)

			tokenValidator := mock.Mock[TokenValidator]()

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{
				Type: PUBLIC,
			}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, alwaysEnabled)
			result, err := interceptor(next)(ctx, req)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(res))
		})

		It("should invoke the next handler when the enabled is false", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(false, nil)

			tokenValidator := mock.Mock[TokenValidator]()
			mock.When(tokenValidator.Validate(mock.AnyContext(), mock.AnyString())).ThenReturn(claims, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{
				Type: SINGLE,
			}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, func() bool {
				return false
			})
			result, err := interceptor(next)(ctx, req)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(res))
		})

		It("should return a permission denied error when the check returns false", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(false, nil)

			tokenValidator := mock.Mock[TokenValidator]()
			mock.When(tokenValidator.Validate(mock.AnyContext(), mock.AnyString())).ThenReturn(claims, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, alwaysEnabled)
			result, err := interceptor(next)(ctx, req)
			Expect(err.Error()).To(Equal("permission_denied: permission denied"))
			Expect(result).To(BeNil())
		})

		It("should return a permission denied error when the request is unauthenticated", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(false, nil)

			extractor := func(req connect.AnyRequest) (string, error) {
				return "", fmt.Errorf("unauthenticated")
			}

			tokenValidator := mock.Mock[TokenValidator]()
			mock.When(tokenValidator.Validate(mock.AnyContext(), mock.AnyString())).ThenReturn(claims, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, extractor, tokenValidator, claimsMapper, alwaysEnabled)
			result, err := interceptor(next)(ctx, req)
			Expect(err.Error()).To(Equal("permission_denied: permission denied"))
			Expect(result).To(BeNil())
		})

		It("should return the error when the check call fails", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(false, fmt.Errorf("unknown error"))

			tokenValidator := mock.Mock[TokenValidator]()
			mock.When(tokenValidator.Validate(mock.AnyContext(), mock.AnyString())).ThenReturn(claims, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, alwaysEnabled)
			result, err := interceptor(next)(ctx, req)
			Expect(err.Error()).To(Equal("unknown error"))
			Expect(result).To(BeNil())
		})
	})

	When("the request is not Checkable", func() {
		It("should return a permission denied error", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(true, nil)

			tokenValidator := mock.Mock[TokenValidator]()

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn("")
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, alwaysEnabled)
			result, err := interceptor(next)(ctx, req)
			Expect(err.Error()).To(Equal("permission_denied: permission denied"))
			Expect(result).To(BeNil())
		})
	})
})
