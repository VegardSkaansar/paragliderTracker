package paragliderdb

import (
	"gopkg.in/mgo.v2/bson"
)

// Webhookinfo handles a webhook request
type Webhookinfo struct {
	ID                 bson.ObjectId `bson:"_id,omitempty" json:"-"`
	URL                string        `json:"webhookURLL"`
	MinTriggerValue    int           `json:"minTriggerValue" bson:"min_trigger_value"`
	LatestInvokedTrack bson.ObjectId `bson:"latest_invoked_track" json:"-"`
	WebhookID          string        `bson:"webhookid" json:"webhookid"`
}

// WebhookResponse gives a respod when get request
type WebhookResponse struct {
	TLatest    int64     `json:"t_latest"`
	Tracks     []TrackID `json:"tracks"`
	Processing int64     `json:"processing"`
	URL        string
}
