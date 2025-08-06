package analytics

import (
	"time"

	"github.com/google/uuid"
)

type EventSchema struct {
	Name      string `json:"name" bson:"name"`
	ProjectId string `json:"project_id"`
	// SessionId uuid.UUID `json:"sessionId"`
	// UserId    uuid.UUID `json:"userId" bson:"userId"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Metadata  Metadata  `json:"metadata" bson:"metadata"`
}

type Metadata struct {
	Page     string `json:"page" bson:"page"`
	ButtonId string `json:"button_id" bson:"button_id"`
	Env      string `json:"env" bson:"env"`
}
type UserSchema struct {
	Id         uuid.UUID      `json:"id"`
	ExternalId string         `json:"external_id"`
	CreatedAt  time.Time      `json:"created_at"`
	Email      string         `json:"email"`
	Name       string         `json:"name"`
	Properties map[string]any `json:"properties"`
}
type ProjectSchema struct {
	Id       uuid.UUID      `json:"id"`
	Url      string         `json:"url"`
	DataTags map[string]int `json:"data_tags"`
}
