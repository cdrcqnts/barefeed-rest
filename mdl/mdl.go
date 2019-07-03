package mdl

import (
	"time"
)

// Link provides the schema by which URLs are saved in the database
type Link struct {
	SID string `json:"sid" bson:"sid" binding:"required"` // Slot ID
	CID string `json:"cid" bson:"cid" binding:"required"` // Channel ID
	URL string `json:"url" bson:"url" binding:"required"` // URL of podcast channel
}

// Channel provides the schema by which feed content is delivered to the client
// TODO: rename to Feed
type Channel struct {
	Updated     time.Time `json:"updated"`
	SID         string    `json:"sid"`
	CID         string    `json:"cid"`
	Url         string    `json:"url"`
	Web         string    `json:"web"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Podcasts    []Podcast `json:"podcasts"`
}

// Podcast is contained by Channel
type Podcast struct {
	Size        int       `json:"size"`
	Released    time.Time `json:"released"`
	PID         string    `json:"pid"`
	Url         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    string    `json:"duration"`
	Image       string    `json:"image"`
}
