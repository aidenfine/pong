package analytics

import (
	"database/sql"

	"github.com/aidenfine/pong/internal/handler/analytics"
	"github.com/go-chi/chi/v5"
)

func Routes(client *sql.DB) chi.Router {
	r := chi.NewRouter()
	handler := &analytics.AnalyticsHandler{DB: client}
	r.Post("/event", handler.PostEvent(client))
	return r
}
