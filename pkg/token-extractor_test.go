package connectpermify

import (
	"connectrpc.com/connect"
	"fmt"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
)

func TestExtract(t *testing.T) {
	t.Run("returns an unauthenticated error when the authorization header is empty", func(t *testing.T) {
		mock.SetUp(t)
		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Header()).ThenReturn(map[string][]string{})

		_, err := DefaultTokenExtractor(req)
		assert.ErrorContains(t, err, "unauthenticated")
	})

	t.Run("returns the token when the authorization header is present and the token type is bearer", func(t *testing.T) {
		mock.SetUp(t)
		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Header()).ThenReturn(http.Header{
			"Authorization": {
				fmt.Sprintf("bearer %s", token),
			},
		})

		result, err := DefaultTokenExtractor(req)
		assert.NoError(t, err)
		assert.Equal(t, token, result)
	})

	t.Run("returns an unauthenticated error when the authorization header is present but the token type is not bearer", func(t *testing.T) {
		mock.SetUp(t)
		req := mock.Mock[connect.AnyRequest]()
		mock.When(req.Header()).ThenReturn(http.Header{
			"Authorization": {
				fmt.Sprintf("bearerer %s", token),
			},
		})

		_, err := DefaultTokenExtractor(req)
		assert.ErrorContains(t, err, "unauthenticated")
	})
}
