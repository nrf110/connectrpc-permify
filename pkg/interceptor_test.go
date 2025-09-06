package connectpermify

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
)

func TestInterceptor_Checkable(t *testing.T) {
	t.Run("invokes the next handler when the check call returns true", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()
		mock.When(client.Check(mock.AnyContext(), mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(true, nil)

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})

		authenticator := mock.Mock[Authenticator]()
		mock.When(authenticator.Authenticate(mock.AnyContext(), mock.Any[connect.AnyRequest]())).ThenReturn(
			&AuthenticationResult{
				Principal: &Resource{
					ID:   "1234",
					Type: "User",
				},
				Attributes: Attributes{},
			},
			nil,
		)

		interceptor := NewPermifyInterceptor(client, authenticator, alwaysEnabled)
		result, err := interceptor(next)(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, res, result)
	})

	t.Run("invokes the next handler when the CheckConfig is public", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()
		mock.When(client.Check(mock.AnyContext(), mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(true, nil)

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{
			IsPublic: true,
		}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})

		authenticator := mock.Mock[Authenticator]()

		interceptor := NewPermifyInterceptor(client, authenticator, alwaysEnabled)
		result, err := interceptor(next)(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, res, result)
	})

	t.Run("invokes the next handler when the enabled flag is false", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()
		mock.When(client.Check(mock.AnyContext(), mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(false, nil)

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})

		authenticator := mock.Mock[Authenticator]()

		interceptor := NewPermifyInterceptor(client, authenticator, func() bool {
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
		mock.When(client.Check(mock.AnyContext(), mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(false, nil)

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})

		authenticator := mock.Mock[Authenticator]()
		mock.When(authenticator.Authenticate(mock.AnyContext(), mock.Any[connect.AnyRequest]())).ThenReturn(
			&AuthenticationResult{
				Principal: &Resource{
					ID:   "1234",
					Type: "User",
				},
				Attributes: Attributes{},
			},
			nil,
		)

		interceptor := NewPermifyInterceptor(client, authenticator, alwaysEnabled)
		result, err := interceptor(next)(ctx, req)
		assert.ErrorContains(t, err, "permission_denied: permission denied")
		assert.Nil(t, result)
	})

	t.Run("returns a permission denied error when the request is unauthenticated", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})

		authenticator := mock.Mock[Authenticator]()
		mock.When(authenticator.Authenticate(mock.AnyContext(), mock.Any[connect.AnyRequest]())).ThenReturn(
			nil,
			connect.NewError(connect.CodePermissionDenied, errors.New("permission denied")),
		)

		interceptor := NewPermifyInterceptor(client, authenticator, alwaysEnabled)
		result, err := interceptor(next)(ctx, req)
		assert.ErrorContains(t, err, "permission_denied: permission denied")
		assert.Nil(t, result)
	})

	t.Run("returns the error when the check call fails", func(t *testing.T) {
		mock.SetUp(t)
		ctx := context.Background()
		client := mock.Mock[CheckClient]()
		expectedErr := fmt.Errorf("unknown error")
		mock.When(client.Check(mock.AnyContext(), mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(false, expectedErr)

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})

		authenticator := mock.Mock[Authenticator]()
		mock.When(authenticator.Authenticate(mock.AnyContext(), mock.Any[connect.AnyRequest]())).ThenReturn(
			&AuthenticationResult{
				Principal: &Resource{
					ID:   "1234",
					Type: "User",
				},
				Attributes: Attributes{},
			},
			nil,
		)

		interceptor := NewPermifyInterceptor(client, authenticator, alwaysEnabled)
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
		mock.When(client.Check(mock.AnyContext(), mock.Any[*Resource](), mock.Any[Attributes](), mock.Any[CheckConfig]())).
			ThenReturn(true, nil)

		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Any()).ThenReturn("")
		res := mock.Mock[connect.AnyResponse]()
		next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			return res, nil
		})

		authenticator := mock.Mock[Authenticator]()

		interceptor := NewPermifyInterceptor(client, authenticator, alwaysEnabled)
		result, err := interceptor(next)(ctx, req)
		assert.ErrorContains(t, err, "permission_denied: permission denied")
		assert.Nil(t, result)
	})
}
