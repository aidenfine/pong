package status

import (
	"time"
)

type CreateStatusUpdateBody struct {
	Message string `json:"message"`
	Service string `json:"service"`
	Status  string `json:"status"`
}

type SnapshotBody struct {
	UserId           string    `json:"userId" bson:"userId"`
	Service          string    `json:"service" bson:"service"`
	CreatedAt        time.Time `json:"createdAt" bson:"createdAt"`
	TotalDataPoints  int       `json:"totalDataPoints" bson:"totalDataPoints"`
	DownDataPoints   int       `json:"downDataPoints" bson:"downDataPoints"`
	UptimePercentage float64   `json:"uptimePercentage" bson:"uptimePercentage"`
}

type GetSnapshotBody struct {
	UserId string `json:"userId"`
	// add more models like query for specific time frames
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
