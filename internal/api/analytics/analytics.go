package analytics

import (
	"github.com/aidenfine/pong/internal/handler/analytics"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Routes(client *mongo.Client) chi.Router {
	r := chi.NewRouter()
	handler := &analytics.AnalyticsHandler{DB: client}
	r.Post("/event", handler.PostEvent(client))
	return r
}
