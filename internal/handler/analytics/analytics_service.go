package analytics

import (
	"context"
	"database/sql"
	"net/http"
	"os"

	common "github.com/aidenfine/pong/internal/handler/common"
	model "github.com/aidenfine/pong/internal/models/analytics"
	"go.uber.org/zap"
)

type AnalyticsHandler struct {
	DB *sql.DB
}

var logger, _ = zap.NewProduction()

func (h *AnalyticsHandler) PostEvent(client *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := common.DecodeJSONBody[model.EventSchema](w, r)
		if err != nil {
			common.Error(w, http.StatusBadRequest, "Invalid Json")
			return
		}
		coll := h.DB.Database(getDatabaseName()).Collection(getAnalyticsCollectionName())
		_, err = coll.InsertOne(context.TODO(), event)
		if err != nil {
			logger.Error("Insert failed", zap.Error(err))
			common.Error(w, http.StatusInternalServerError, "Something went wrong")
			return
		}
		logger.Info("Created event")
		common.WriteJSON(w, http.StatusCreated, event)
	}
}

func getDatabaseName() string {
	return os.Getenv("GO_ENV")
}
func getAnalyticsCollectionName() string {
	return os.Getenv("ANALYTICS_COLLECTION_NAME")
}
