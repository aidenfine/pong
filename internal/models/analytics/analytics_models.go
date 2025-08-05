package analytics

import (
	"time"

	"github.com/google/uuid"
)

type EventSchema struct {
	Name      string    `json:"name" bson:"name"`
	SessionId uuid.UUID `json:"sessionId"`
	// UserId    uuid.UUID `json:"userId" bson:"userId"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Metadata  Metadata  `json:"metadata" bson:"metadata"`
}

type Metadata struct {
	Page     string `json:"page" bson:"page"`
	ButtonId string `json:"buttonId" bson:"buttonId"`
	Env      string `json:"env" bson:"env"`
}
type UserSchema struct {
	Id         uuid.UUID      `json:"id"`
	ExternalId string         `json:"externalId"`
	CreatedAt  time.Time      `json:"createdAt"`
	Email      string         `json:"email"`
	Name       string         `json:"name"`
	Properties map[string]any `json:"properties"`
}
