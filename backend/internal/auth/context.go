package auth

import "context"

type contextKey string

const claimsKey contextKey = "claims"

// SetClaimsInContext stores JWT claims in the request context.
func SetClaimsInContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// GetClaimsFromContext retrieves JWT claims from the request context.
func GetClaimsFromContext(ctx context.Context) *Claims {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	if !ok {
		return nil
	}
	return claims
}
