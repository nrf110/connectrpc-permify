package connectpermify

import "context"

type permifySnapToken struct{}

var PermifySnapToken permifySnapToken = permifySnapToken{}

func GetSnapToken(ctx context.Context) string {
	value := ctx.Value(PermifySnapToken)
	if value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func SetSnapToken(ctx context.Context, snapToken string) context.Context {
	return context.WithValue(ctx, PermifySnapToken, snapToken)
}
