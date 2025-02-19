package connectpermify

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterceptor_Checkable(t *testing.T) {
	claims := &validator.ValidatedClaims{
		RegisteredClaims: validator.RegisteredClaims{
			Subject: "abcde",
		},
	}

	t.Run("invokes the next handler when the check call returns true", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()
		mock.When(client.Check(mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(true, nil)

		tokenValidator := mock.Mock[TokenValidator]()
		mock.When(tokenValidator.Validate(mock.AnyContext(), mock.AnyString())).
			ThenReturn(claims, nil)

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})
		interceptor := NewPermifyInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, alwaysEnabled)
		result, err := interceptor(next)(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, res, result)
	})

	t.Run("invokes the next handler when the CheckConfig is public", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()
		mock.When(client.Check(mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(true, nil)

		tokenValidator := mock.Mock[TokenValidator]()

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{
			Type: PUBLIC,
		}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})
		interceptor := NewPermifyInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, alwaysEnabled)
		result, err := interceptor(next)(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, res, result)
	})

	t.Run("invokes the next handler when the enabled flag is false", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()
		mock.When(client.Check(mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(false, nil)

		tokenValidator := mock.Mock[TokenValidator]()
		mock.When(tokenValidator.Validate(mock.AnyContext(), mock.AnyString())).
			ThenReturn(claims, nil)

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{
			Type: SINGLE,
		}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})
		interceptor := NewPermifyInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, func() bool {
			return false
		})
		result, err := interceptor(next)(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, res, result)
	})

	t.Run("returns a permission denied error when the check returns false", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()
		mock.When(client.Check(mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(false, nil)

		tokenValidator := mock.Mock[TokenValidator]()
		mock.When(tokenValidator.Validate(mock.AnyContext(), mock.AnyString())).
			ThenReturn(claims, nil)

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})
		interceptor := NewPermifyInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, alwaysEnabled)
		result, err := interceptor(next)(ctx, req)
		assert.ErrorContains(t, err, "permission_denied: permission denied")
		assert.Nil(t, result)
	})

	t.Run("returns a permission denied error when the request is unauthenticated", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()
		mock.When(client.Check(mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(false, nil)

		extractor := func(req connect.AnyRequest) (string, error) {
			return "", fmt.Errorf("unauthenticated")
		}

		tokenValidator := mock.Mock[TokenValidator]()
		mock.When(tokenValidator.Validate(mock.AnyContext(), mock.AnyString())).
			ThenReturn(claims, nil)

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})
		interceptor := NewPermifyInterceptor(client, extractor, tokenValidator, claimsMapper, alwaysEnabled)
		result, err := interceptor(next)(ctx, req)
		assert.ErrorContains(t, err, "permission_denied: permission denied")
		assert.Nil(t, result)
	})

	t.Run("returns the error when the check call fails", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()
		expectedErr := fmt.Errorf("unknown error")
		mock.When(client.Check(mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(false, expectedErr)

		tokenValidator := mock.Mock[TokenValidator]()
		mock.When(tokenValidator.Validate(mock.AnyContext(), mock.AnyString())).
			ThenReturn(claims, nil)

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})
		interceptor := NewPermifyInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, alwaysEnabled)
		result, err := interceptor(next)(ctx, req)

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, result)
	})
}

func TestInterceptor_NotCheckable(t *testing.T) {
	t.Run("returns a permission denied error", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()
		mock.When(client.Check(mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(true, nil)

		tokenValidator := mock.Mock[TokenValidator]()

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn("")
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})
		interceptor := NewPermifyInterceptor(client, tokenExtractor, tokenValidator, claimsMapper, alwaysEnabled)
		result, err := interceptor(next)(ctx, req)
		assert.ErrorContains(t, err, "permission_denied: permission denied")
		assert.Nil(t, result)
	})
}
