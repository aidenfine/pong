package status

import (
	"time"
)

type CreateStatusUpdateBody struct {
	Message string `json:"message" bson:"message"`
	Service string `json:"service" bson:"service"`
	Status  string `json:"status" bson:"status"`
}

type SnapshotBody struct {
	Service          string    `json:"service" bson:"service"`
	Timestamp        time.Time `json:"timestamp" bson:"timestamp"`
	TotalDataPoints  int       `json:"totalDataPoints" bson:"totalDataPoints"`
	DownDataPoints   int       `json:"downDataPoints" bson:"downDataPoints"`
	UptimePercentage float64   `json:"uptimePercentage" bson:"uptimePercentage"`
}

type GetStatusBody struct {
	Service string `json:"service" bson:"service"`
}

type StatusUpdateBody struct {
	CreateStatusUpdateBody
}

type LiveStatusModel struct {
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Service   string    `json:"service" bson:"service"`
	Status    string    `json:"status" bson:"status"`
}
