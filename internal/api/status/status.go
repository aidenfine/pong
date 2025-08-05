package status

import (
	"github.com/aidenfine/pong/internal/handler/status"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Routes(client *mongo.Client) chi.Router {
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
