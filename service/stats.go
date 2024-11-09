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

	impressions := storage.GetImpressionsByCampaign(campaignID)

	stats := &model.Stats{
		CampaignID: campaignID,
	}

	for _, impression := range impressions {
		if time.Since(impression.Timestamp) <= time.Hour {
			stats.LastHour++
		}
		if time.Since(impression.Timestamp) <= 24*time.Hour {
			stats.LastDay++
		}
		stats.TotalCount++
	}

	return stats, nil
}
