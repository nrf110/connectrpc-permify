package connectpermify

import (
	"fmt"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-viper/mapstructure/v2"
)

type ClaimsMapper func(claims *validator.ValidatedClaims) (*Resource, Attributes, error)

type CustomClaimsMapper[T validator.CustomClaims] func(T) (Attributes, error)

func DefaultCustomClaimsMapper[T validator.CustomClaims]() CustomClaimsMapper[T] {
	return func(t T) (Attributes, error) {
		attributes := Attributes{}
		err := mapstructure.Decode(t, &attributes)
		if err != nil {
			return nil, err
		}
		return attributes, nil
	}
}

func DefaultClaimsMapper[T validator.CustomClaims](customClaimsMapper CustomClaimsMapper[T]) ClaimsMapper {
	return func(claims *validator.ValidatedClaims) (*Resource, Attributes, error) {
		subject := claims.RegisteredClaims.Subject
		customClaims, ok := claims.CustomClaims.(T)
		if !ok {
			return nil, nil, fmt.Errorf("unexpected custom claims type")
		}

		attributes, err := customClaimsMapper(customClaims)
		if err != nil {
			return nil, nil, err
		}

		return &Resource{
			ID:   subject,
			Type: "user",
		}, attributes, nil
	}
}
