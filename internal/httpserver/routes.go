package httpserver

import (
	"net/http"

	authMiddleware "github.com/caneel7/keynote-api/internal/common/middleware/auth"
	"github.com/caneel7/keynote-api/internal/common/tokens"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/features/auth"
	"github.com/caneel7/keynote-api/internal/features/project"
	"github.com/caneel7/keynote-api/internal/features/tenant"
	sqlc "github.com/caneel7/keynote-api/internal/db/sqlc"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)



func newApi (router chi.Router) huma.API {
	return humachi.New(
		router,
		huma.DefaultConfig("keynote-api","0.1.0"),
	)
}

func routes (
	db *gen.Client,
	sqlc *sqlc.Queries,
	signer *tokens.Signer,
) http.Handler {

	router := chi.NewRouter()
	api := newApi(router)

	protectedAPI := huma.NewGroup(api)
	protectedAPI.UseMiddleware(authMiddleware.VerifyJWT(api,signer))

	healthHandler(api)

	/*
	* auth api mounted : ["/login","/register","/me","/refresh","/.well-known/jwks.json"]
	*/
	auth.Mount(api,protectedAPI,db,signer)

	/*
	* tenat api mounted : ["/tenant"]
	*/
	tenant.Mount(protectedAPI,db)

	/*
	* project api mounted : ["/project"]
	*/
	project.Mount(protectedAPI,db)

	return router;
}