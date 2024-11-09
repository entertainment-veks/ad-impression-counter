package service

import (
	"ad_impression_counter/model"
	"ad_impression_counter/storage"
	"log"
	"time"
)

func GetCampaignStats(campaignID string) (*model.Stats, error) {
	_, err := GetCampaign(campaignID)
	if err != nil {
		log.Printf("failed to get campaign by ID: %s: %v", campaignID, err)
		return nil, err
	}

	impressions, err := storage.GetImpressionsByCampaign(campaignID)
	if err != nil {
		log.Printf("failed to get impressions by campaign ID: %s: %v", campaignID, err)
		return nil, err
	}

	stats := &model.Stats{
		CampaignID: campaignID,
		LastHour:   0,
		LastDay:    0,
		TotalCount: int64(len(impressions)),
	}

	// Since impressions are sorted by timestamp, we can iterate them in reverse order to not check impressions older than 24 hours
	for i := len(impressions) - 1; i >= 0; i-- {
		if time.Since(impressions[i].Timestamp) < time.Hour {
			stats.LastHour++
		}

		if time.Since(impressions[i].Timestamp) < 24*time.Hour {
			stats.LastDay++
		}

		if time.Since(impressions[i].Timestamp) >= 24*time.Hour {
			break // No need to check older impressions, because total count is already calculated by len(impressions)
		}
	}

	return stats, nil
}
