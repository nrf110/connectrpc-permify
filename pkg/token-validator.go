package connectpermify

import (
	"context"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"net/url"
	"time"
)

type TokenValidator interface {
	Validate(ctx context.Context, token string) (*validator.ValidatedClaims, error)
}

type OIDCConfig struct {
	TrustedIssuer      string
	Audiences          []string
	SignatureAlgorithm validator.SignatureAlgorithm
	CustomClaims       func() validator.CustomClaims
}

type OIDCTokenValidator struct {
	validator *validator.Validator
}

func NewOIDCTokenValidator(config OIDCConfig) (*OIDCTokenValidator, error) {
	uri, err := url.Parse(config.TrustedIssuer)
	if err != nil {
		return nil, err
	}
	provider := jwks.NewCachingProvider(uri, 1*time.Hour)
	v, err := validator.New(
		provider.KeyFunc,
		config.SignatureAlgorithm,
		config.TrustedIssuer,
		config.Audiences,
		validator.WithCustomClaims(config.CustomClaims))

	return &OIDCTokenValidator{
		validator: v,
	}, nil
}

func (authn *OIDCTokenValidator) Validate(ctx context.Context, token string) (*validator.ValidatedClaims, error) {
	claims, err := authn.validator.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return claims.(*validator.ValidatedClaims), nil
}
