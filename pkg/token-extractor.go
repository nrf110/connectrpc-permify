package connectpermit

import (
	"connectrpc.com/connect"
	"fmt"
	"strings"
)

type TokenExtractor func(req connect.AnyRequest) (string, error)

func DefaultTokenExtractor(req connect.AnyRequest) (string, error) {
	header := req.Header().Get("Authorization")
	if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return "", fmt.Errorf("unauthenticated")
	}
	return header[7:], nil
}
