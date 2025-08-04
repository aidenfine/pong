package analytics

import "time"

type EventSchema struct {
	Event     string    `json:"event" bson:"event"`
	UserId    string    `json:"userId" bson:"userId"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Metadata  Metadata  `json:"metadata" bson:"metadata"`
}

type Metadata struct {
	Page     string `json:"page" bson:"page"`
	ButtonId string `json:"buttonId" bson:"buttonId"`
	Env      string `json:"env" bson:"env"`
}
