package connectpermify

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

func NewPermifyInterceptor(
	client CheckClient,
	tokenExtractor TokenExtractor,
	tokenValidator TokenValidator,
	claimsMapper ClaimsMapper,
	enabled func() bool,
) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			checkable, ok := req.Any().(Checkable)
			if !ok {
				return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
			}
			checks := checkable.GetChecks()
			if enabled() && !checks.IsPublic() {
				token, err := tokenExtractor(req)
				if err != nil {
					return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
				}

				claims, err := tokenValidator.Validate(ctx, token)
				if err != nil {
					return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
				}

				principal, attributes, err := claimsMapper(claims)
				if err != nil {
					return nil, err
				}

				result, err := client.Check(principal, attributes, checks)
				if err != nil {
					return nil, err
				}
				if !result {
					return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
				}
			}
			return next(ctx, req)
		})
	}
}
