package status

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	common "github.com/aidenfine/pong/internal/handler/common"
	model "github.com/aidenfine/pong/internal/models/status"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

var POLLING_RATE = 30
var SNAPSHOT_CREATION_RATE = 2

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
		coll := h.DB.Database(getDatabaseName()).Collection(getCollectionName())
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
		response := map[string]any{
			"message": "Created",
			"data": map[string]any{
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
			ticker := time.NewTicker(time.Duration(POLLING_RATE) * time.Second)
			defer ticker.Stop()
			coll := client.Database(getDatabaseName()).Collection(getLiveStatusCollectionName())
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Stopped polling")
					return
				case <-ticker.C:
					up := checkService(url)
					status := "ERROR"
					if up {
						status = "OK"
					}
					doc := bson.M{
						"service":   url,
						"status":    status,
						"timestamp": time.Now().UTC(),
					}
					_, err := coll.InsertOne(ctx, doc)
					if err != nil {
						logger.Error("Insert Failed", zap.Error(err))
					} else {
						fmt.Println("Inserted doc", doc)
					}
				}

			}
		}(url)
	}
}

func (h *StatusHandler) StartLiveSnapshot(client *mongo.Client) http.HandlerFunc {
	url := os.Getenv("URL")
	ctx := context.Background()

	return func(w http.ResponseWriter, r *http.Request) {
		go func(url string) {
			ticker := time.NewTicker(time.Duration(SNAPSHOT_CREATION_RATE) * time.Minute)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					fmt.Println("Stopped polling")
					return

				case <-ticker.C:
					var timestampToQuery time.Time = time.Now().UTC().AddDate(0, 0, -7)
					latestSnapshot, err := h.fetchLatestSnapshot(getUrlenv())
					if err != nil {
						logger.Error("Failed to fetch latest snapshot", zap.Error(err))
					} else {
						fmt.Println("Found snapshot using", latestSnapshot)
						timestampToQuery = latestSnapshot.Timestamp
					}

					countMap, err := h.fetchDatapointsSinceDate(getUrlenv(), timestampToQuery)
					if err != nil {
						logger.Error("Failed to fetch datapoints", zap.Error(err))
						continue
					}

					totalCountSinceLastSnapshot := 0
					for _, val := range countMap {
						totalCountSinceLastSnapshot += val
					}

					if totalCountSinceLastSnapshot == 0 {
						logger.Warn("No new datapoints since last snapshot")
						continue
					}

					newTotalDown := latestSnapshot.DownDataPoints + countMap["ERROR"]
					newTotalDataPoints := latestSnapshot.TotalDataPoints + totalCountSinceLastSnapshot

					uptimePercentage := calculatePercentage(float64(countMap["ERROR"]), float64(totalCountSinceLastSnapshot))

					_, err = client.Database(getDatabaseName()).Collection(getSnapshotCollectionName()).InsertOne(context.TODO(), bson.M{
						"service":          getUrlenv(),
						"timestamp":        time.Now().UTC(),
						"totalDownPoints":  newTotalDown,
						"downDataPoints":   countMap["ERROR"],
						"totalDataPoints":  newTotalDataPoints,
						"uptimePercentage": uptimePercentage,
					})
					if err != nil {
						logger.Error("Insert failed", zap.Error(err))
						return
					}
				}
			}
		}(url)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Live snapshot started"))
	}
}
func (h *StatusHandler) GetAnalytics(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := common.DecodeJSONBody[model.SnapshotBody](w, r)
		if err != nil {
			common.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		snapshotBody, err := h.fetchLatestSnapshot(reqBody.Service)
		fmt.Println(snapshotBody.Timestamp, "timestamp when snapshot was created")
		if err != nil {
			if err == mongo.ErrNoDocuments {
				fmt.Println("no docs")
				common.Error(w, http.StatusEarlyHints, "Uptime percentage is not avaliable yet or does not exist")
				return
			}
			common.Error(w, http.StatusBadGateway, "An error has occured getting an item from the database")
			return
		}
		common.WriteJSON(w, http.StatusOK, snapshotBody)
	}
}

func (h *StatusHandler) TestInsert(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap, err := common.DecodeJSONBody[model.SnapshotBody](w, r)
		if err != nil {
			common.Error(w, http.StatusBadRequest, "invalid json")
			return
		}

		if snap.Timestamp.IsZero() {
			snap.Timestamp = time.Now().UTC()
		}

		coll := h.DB.Database(getDatabaseName()).Collection(getSnapshotCollectionName())

		doc := bson.M{
			"service":          snap.Service,
			"timestamp":        snap.Timestamp,
			"totalDataPoints":  snap.TotalDataPoints,
			"downDataPoints":   snap.DownDataPoints,
			"uptimePercentage": snap.UptimePercentage,
		}

		_, err = coll.InsertOne(context.TODO(), doc)
		if err != nil {
			logger.Error("Insert failed", zap.Error(err))
			http.Error(w, "insert failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Snapshot inserted successfully"))
	}
}

func checkService(url string) bool {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		return false
	}
	return true
}

func (h *StatusHandler) fetchLatestSnapshot(service string) (model.SnapshotBody, error) {
	var snapshotBody model.SnapshotBody
	coll := h.DB.Database(getDatabaseName()).Collection(getSnapshotCollectionName())
	filter := bson.M{"service": service}
	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	err := coll.FindOne(context.TODO(), filter, opts).Decode(&snapshotBody)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No snapshot found")
			return snapshotBody, err
		}
		logger.Error("Error finding snapshot", zap.Error(err))
		return snapshotBody, err
	}
	return snapshotBody, nil
}
func (h *StatusHandler) fetchDatapointsSinceDate(service string, timestamp time.Time) (map[string]int, error) {
	fmt.Println("checking timestamp", timestamp)
	fmt.Println("checking", service)
	filter := bson.M{
		"service":   service,
		"timestamp": bson.M{"$gte": timestamp},
	}

	cursor, err := h.DB.Database(getDatabaseName()).Collection(getLiveStatusCollectionName()).Find(context.TODO(), filter)
	if err != nil {
		logger.Error("error querying DB", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(context.TODO())

	countMap := make(map[string]int)

	for cursor.Next(context.TODO()) {
		var doc model.LiveStatusModel
		if err := cursor.Decode(&doc); err != nil {
			logger.Error("error decoding", zap.Error(err))
			continue
		}
		fmt.Println(doc, "returned doc")
		switch doc.Status {
		case "OK":
			countMap["OK"]++
		case "ERROR":
			countMap["ERROR"]++
		}
	}

	if err := cursor.Err(); err != nil {
		logger.Error("cursor error", zap.Error(err))
		return nil, err
	}

	return countMap, nil
}

func calculatePercentage(errCount float64, totalCount float64) float64 {
	return math.Round(100 - (errCount/totalCount)*100)
}

func getDatabaseName() string {
	return os.Getenv("GO_ENV")
}
func getCollectionName() string {
	return os.Getenv("status")

}
func getLiveStatusCollectionName() string {
	return os.Getenv("LIVE_STATUS_COLLECTION_NAME")
}
func getSnapshotCollectionName() string {
	return os.Getenv("SNAPSHOT_COLLECTION_NAME")
}

func getUrlenv() string {
	return os.Getenv(("URL"))
}
