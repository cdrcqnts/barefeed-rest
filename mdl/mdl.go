package mdl

import (
	"time"
)

type Feed struct {
	SID string `json:"sid" bson:"sid" binding:"required"` // Slot ID
	CID string `json:"cid" bson:"cid" binding:"required"` // Channel ID
	URL string `json:"url" bson:"url" binding:"required"` // URL of podcast channel
}

type Channel struct {
	SID         string    `json:"sid"`
	CID         string    `json:"cid"`
	Url         string    `json:"url"`
	Web         string    `json:"web"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Updated     time.Time `json:"updated"`
	Image       string    `json:"image"`
	Podcasts    []Podcast `json:"podcasts"`
}

type Podcast struct {
	PID         string    `json:"pid"`
	Url         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    string    `json:"duration"`
	Released    time.Time `json:"released"`
	Image       string    `json:"image"`
	Size        int       `json:"size"`
}
