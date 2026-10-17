package rbac

import (
	"context"
	"net/http"
	"strings"

	"github.com/caneel7/keynote-api/internal/constants"
	"github.com/caneel7/keynote-api/internal/ent/gen"

	"github.com/danielgtaylor/huma/v2"
)

type Authorizer struct {
	ctx context.Context
	db  *gen.Client
	api huma.API
}

func New(ctx context.Context, db *gen.Client, api huma.API) *Authorizer {
	return &Authorizer{ctx, db, api}
}

func (a *Authorizer) Require(key string) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

		keyModule := strings.Split(key, ":")[0]

		err := a.checkModule(keyModule,ctx)

		if err != nil {
			huma.WriteErr(a.api,ctx,http.StatusUnauthorized,"Invalid authorization header")
			return
		}
		next(ctx)
	}
}

func (a *Authorizer) checkModule(key string, ctx huma.Context) (error) {

	module, _ := constants.ResourceToScope[key]

	var headers = map[string]string{}

	for _, m := range(module) {
		headers[m] = ctx.Header(m)
	}

	return nil
}