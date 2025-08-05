package status

import (
	"database/sql"

	"github.com/aidenfine/pong/internal/handler/status"
	"github.com/go-chi/chi/v5"
)

func Routes(client *sql.DB) chi.Router {
	r := chi.NewRouter()
	handler := &status.StatusHandler{DB: client}
	r.Get("/health", status.Health())
	r.Get("/poll", handler.StartPolling(client))
	r.Post("/", handler.CreateStatusUpdate(client))
	r.Post("/test", handler.TestInsert(client))
	r.Get("/analytics", handler.GetAnalytics(client))
	r.Get("/start-snapshots", handler.StartLiveSnapshot(client))
	return r
}
