package connectpermify

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

type OAuth2Authenticator struct {
	tokenExtractor TokenExtractor
	tokenValidator TokenValidator
	claimsMapper   ClaimsMapper
}

type Opt func(*OAuth2Authenticator)

func WithTokenExtractor(tokenExtractor TokenExtractor) Opt {
	return func(oa *OAuth2Authenticator) {
		oa.tokenExtractor = tokenExtractor
	}
}

func WithTokenValidator(tokenValidator TokenValidator) Opt {
	return func(oa *OAuth2Authenticator) {
		oa.tokenValidator = tokenValidator
	}
}

func WithClaimsMapper(claimsMapper ClaimsMapper) Opt {
	return func(oa *OAuth2Authenticator) {
		oa.claimsMapper = claimsMapper
	}
}

func NewOAuth2Authenticator[T validator.CustomClaims](tokenValidator TokenValidator) *OAuth2Authenticator {
	authn := OAuth2Authenticator{
		tokenExtractor: DefaultTokenExtractor,
		tokenValidator: tokenValidator,
		claimsMapper: DefaultClaimsMapper(
			DefaultCustomClaimsMapper[T](),
		),
	}

	return &authn
}

func (oauth OAuth2Authenticator) Authenticate(ctx context.Context, req connect.AnyRequest) (*AuthenticationResult, error) {
	token, err := oauth.tokenExtractor(req)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	claims, err := oauth.tokenValidator.Validate(ctx, token)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	principal, attributes, err := oauth.claimsMapper(claims)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	return &AuthenticationResult{
		Context:    ctx,
		Principal:  principal,
		Attributes: attributes,
	}, nil
}
