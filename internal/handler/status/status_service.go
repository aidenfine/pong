package status

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	common "github.com/aidenfine/pong/internal/handler/common"
	model "github.com/aidenfine/pong/internal/models/status"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

var DATABASE_NAME = os.Getenv("DB_ENV")
var COLLECTION_NAME = os.Getenv("status")
var LIVE_STATUS_COLLECTION_NAME = os.Getenv("LIVE_STATUS_COLLECTION_NAME")
var POLLING_RATE = 30

type StatusHandler struct {
	DB *mongo.Client
}

var logger, _ = zap.NewProduction()

func Health() http.HandlerFunc {
	urlValue := os.Getenv("URL")
	url := urlValue
	status := checkService(url)
	return func(w http.ResponseWriter, r *http.Request) {
		if status {
			common.Ok(w, map[string]string{
				"status": "OK",
			})
		} else {
			common.Error(w, http.StatusInternalServerError, "Service is down")
		}
	}
}
func (h *StatusHandler) CreateStatusUpdate(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := common.DecodeJSONBody[model.CreateStatusUpdateBody](w, r)
		if err != nil {
			common.Error(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		coll := h.DB.Database(DATABASE_NAME).Collection(COLLECTION_NAME)
		doc := model.CreateStatusUpdateBody{
			Message: reqBody.Message,
			Service: reqBody.Service,
			Status:  reqBody.Status,
		}
		res, err := coll.InsertOne(context.TODO(), doc)
		if err != nil {
			logger.Error("InsertOne failed", zap.Error(err))
			common.Error(w, http.StatusInternalServerError, "There was an issue inserting the document")
			return
		}
		response := map[string]interface{}{
			"message": "Created",
			"data": map[string]interface{}{
				"id":      res.InsertedID,
				"message": reqBody.Message,
				"service": reqBody.Service,
				"status":  reqBody.Status,
			},
		}
		common.WriteJSON(w, http.StatusCreated, response)

	}
}

func (h *StatusHandler) StartPolling(client *mongo.Client) http.HandlerFunc {
	urlValue := os.Getenv("URL")
	url := urlValue
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		go func(url string) {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			coll := client.Database(DATABASE_NAME).Collection(LIVE_STATUS_COLLECTION_NAME)
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Stopped polling")
					return
				case <-ticker.C:
					up := checkService(url)
					status := "DOWN"
					if up {
						status = "UP"
					}
					_, err := coll.InsertOne(ctx, bson.M{
						"serviceId": url,
						"status":    status,
						"timestamp": time.Now().UTC(),
					})
					if err != nil {
						logger.Error("Insert Failed", zap.Error(err))
					} else {
						fmt.Println("Inserted doc")
					}

				}

			}
		}(url)
	}
}

func checkService(url string) bool {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		return false
	}
	return true
}
