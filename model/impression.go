package model

import "time"

type Impression struct {
	CampaignID string    `json:"campaign_id"`
	Timestamp  time.Time `json:"timestamp"`
	UserID     string    `json:"user_id"`
	AdID       string    `json:"ad_id"`
}
