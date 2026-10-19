package httpserver

import (
	"net"
	"net/http"

	"github.com/caneel7/keynote-api/config"
	"github.com/caneel7/keynote-api/internal/common/tokens"
	sqlc "github.com/caneel7/keynote-api/internal/db/sqlc"
	"github.com/caneel7/keynote-api/internal/ent/gen"
)

func NewServer(
	c config.HTTP,
	db *gen.Client,
	sqlc *sqlc.Queries,
	singer *tokens.Signer,
) *http.Server {
	return &http.Server{
		Addr:         net.JoinHostPort("", c.Addr),
		Handler:      routes(db, sqlc, singer),
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	}
}
