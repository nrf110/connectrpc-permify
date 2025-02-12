package connectpermit

import (
	"connectrpc.com/connect"
	"context"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

type testCustomClaims struct {
	validator.CustomClaims
	Roles          []string `json:"roles"`
	OrganizationID string   `json:"organizationId"`
}

func (tcc *testCustomClaims) Validate(ctx context.Context) error {
	return nil
}

type stubCheckable struct {
	Checkable
	checks CheckConfig
}

func (r *stubCheckable) GetChecks() CheckConfig {
	return r.checks
}

func alwaysEnabled() bool { return true }

func tokenExtractor(req connect.AnyRequest) (string, error) {
	return "mocktoken", nil
}

func claimsMapper(_ *validator.ValidatedClaims) (*Resource, Attributes, error) {
	return &Resource{ID: "abcde"}, nil, nil
}
