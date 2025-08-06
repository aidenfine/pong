package analytics

import (
	"database/sql"
	"encoding/json"
	"errors"
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

		query := `INSERT INTO events (id, project_id, name, timestamp, metadata)
          VALUES ($1, $2, $3, $4, $5)`

		_, err = client.Exec(query, id, event.ProjectId, event.Name, event.Timestamp, string(metadataJSON))

		if err != nil {
			logger.Error("Insert failed", zap.Error(err))
			common.Error(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		logger.Info("Event has been created", zap.String("id", id.String()))
		common.WriteJSON(w, http.StatusCreated, event)
	}
}

func (h *AnalyticsHandler) GetProject(client *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("projectId")
		if id == "" {
			common.Error(w, http.StatusBadRequest, "Project id must be sent")
			return
		}

		project, err := h.FetchProjectByID(client, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				common.Error(w, http.StatusNotFound, "Item not found")
				return
			}
			logger.Error("Error getting item from db", zap.Error(err))
			common.Error(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		common.WriteJSON(w, http.StatusOK, project)
	}
}

func (h *AnalyticsHandler) CreateProject(client *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project, err := common.DecodeJSONBody[model.ProjectSchema](w, r)
		if err != nil {
			logger.Error("Failed to decode json", zap.Error(err))
			common.Error(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		dataTagsJson, err := json.Marshal(project.DataTags)
		if err != nil {
			logger.Error("Failed to marshal metadata", zap.Error(err))
			return
		}
		query := `INSERT INTO projects (id, url, data_tags) VALUES ($1, $2, $3)`
		_, err = client.Exec(query, project.Id, project.Url, string(dataTagsJson))

		if err != nil {
			logger.Error("Insert failed", zap.Error(err))
			common.Error(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		logger.Info("Project has been created", zap.String("id", project.Id.String()))
		common.WriteJSON(w, http.StatusCreated, project)
	}
}

// func updateDataTags(db *sql.DB, projectID string, key string, value string) error {
// 	query := `
//         UPDATE projects
//         SET data_tags = json_set(data_tags, ?, ?)
//         WHERE id = ?;
//     `
// 	// json_set path argument needs to be like '$.key'
// 	jsonPath := fmt.Sprintf("$.%s", key)

// 	_, err := db.Exec(query, jsonPath, value, projectID)
// 	return err
// }

// kinda hacky should make this cleaner
func (h *AnalyticsHandler) FetchProjectByID(client *sql.DB, id string) (model.ProjectSchema, error) {
	query := `SELECT id, url, data_tags FROM projects WHERE id = ? LIMIT 1`

	var (
		project     model.ProjectSchema
		rawDataTags map[string]any
	)

	err := client.QueryRow(query, id).Scan(&project.Id, &project.Url, &rawDataTags)
	if err != nil {
		return project, err
	}

	project.DataTags = make(map[string]int)
	for k, v := range rawDataTags {
		if f, ok := v.(float64); ok {
			project.DataTags[k] = int(f)
		} else {
			logger.Warn("Unexpected type in data_tags", zap.Any("key", k), zap.Any("value", v))
		}
	}

	return project, nil
}
