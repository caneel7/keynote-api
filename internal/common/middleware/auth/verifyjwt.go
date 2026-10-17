package auth

import (
	"context"
	"net/http"
	"strings"

	tokens "github.com/caneel7/keynote-api/internal/common/tokens"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type contextKey string

const ClaimKey contextKey = "jwt_claims"

func ClaimFromContext(ctx context.Context) (tokens.Claims,bool) {
	claims, ok := ctx.Value(ClaimKey).(tokens.Claims)
	return claims,ok
}

func GetUserIDFromClaims(ctx context.Context) (uuid.UUID,error) {
	c, _ := ClaimFromContext(ctx)
	p, err := uuid.Parse(c.Subject)
	return p, err
}


func VerifyJWT(api huma.API, signer *tokens.Signer) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context,next func(huma.Context)){
		authHeader := ctx.Header("Authorization")
		if authHeader == "" {
			huma.WriteErr(api,ctx,http.StatusUnauthorized,"Missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader," ",2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			huma.WriteErr(api,ctx,http.StatusUnauthorized,"Invalid authorization header")
			return
		}

		claims, err := signer.Verify(parts[1])
		if err != nil {
			huma.WriteErr(api,ctx,http.StatusUnauthorized,err.Error())
			return
		}
		

		ctx = huma.WithValue(ctx,ClaimKey,claims)
		next(ctx)
	}
}