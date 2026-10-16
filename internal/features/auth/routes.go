package auth

import (
	"github.com/danielgtaylor/huma/v2"

	"github.com/caneel7/keynote-api/internal/common/tokens"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/features/auth/jwks"
	"github.com/caneel7/keynote-api/internal/features/auth/login"
	"github.com/caneel7/keynote-api/internal/features/auth/me"
	"github.com/caneel7/keynote-api/internal/features/auth/refresh"
	"github.com/caneel7/keynote-api/internal/features/auth/register"
	"github.com/caneel7/keynote-api/internal/features/auth/sessions"
)


func Mount(
	public huma.API,
	protected huma.API,
	db *gen.Client,
	signer *tokens.Signer,
) {

	sess := sessions.NewService(signer)

	/*
	* @path : "/.well-known/jwks.json"
	*/
	jwks.JwksRoute(public, signer)

	/*
	* @path : "/register"
	*/
	register.RegisterUserRoute(public, register.NewService(db, sess))

	/*
	* @path : "/login"
	*/
	login.LoginRoute(public, login.NewService(db, sess))

	/*
	* @path : "/refresh"
	*/
	refresh.RefreshTokenRoute(public, refresh.NewService(db, sess))

	/*
	* @path : "/me"
	*/
	me.GetMeRoute(protected, me.NewService(db))

}
