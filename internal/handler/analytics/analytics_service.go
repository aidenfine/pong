package analytics

import (
	"database/sql"
	"encoding/json"
	"net/http"

	common "github.com/aidenfine/pong/internal/handler/common"
	model "github.com/aidenfine/pong/internal/models/analytics"
	"github.com/google/uuid"
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
			logger.Error("Failed to decode json", zap.Error(err))
			common.Error(w, http.StatusBadRequest, "Invalid JSON")

			return
		}
		metadataJSON, err := json.Marshal(event.Metadata)
		if err != nil {
			logger.Error("Failed to marshal metadata", zap.Error(err))
			return
		}

		id := uuid.New()

		query := `INSERT INTO events (id, name, timestamp, metadata)
          VALUES ($1, $2, $3, $4)`

		_, err = client.Exec(query, id, event.Name, event.Timestamp, string(metadataJSON))

		if err != nil {
			logger.Error("Insert failed", zap.Error(err))
			common.Error(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		logger.Info("Event has been created", zap.String("id", id.String()))
		common.WriteJSON(w, http.StatusCreated, event)
	}
}
