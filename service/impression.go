package service

import (
	"ad_impression_counter/config"
	"ad_impression_counter/model"
	"ad_impression_counter/storage"
	"log"
	"time"
)

func TrackImpression(impression model.Impression, cfg config.Config) error {
	campaign, err := GetCampaign(impression.CampaignID)
	if err != nil {
		log.Printf("failed to get campaign by ID: %s: %v", impression.CampaignID, err)
		return err
	}

	if campaign.StartTime.After(impression.Timestamp) {
		log.Printf("campaign has not started yet. starting time: %s", campaign.StartTime)
		return ErrCampaignNotStarted
	}

	oldImpressions, err := storage.GetImpressionsByCampaign(impression.CampaignID)
	if err != nil {
		log.Printf("failed to get impressions by campaign ID: %s: %v", impression.CampaignID, err)
		return err
	}

	// to find the latest impression we have to iterate in reverse order
	for i := len(oldImpressions) - 1; i >= 0; i-- {
		if oldImpressions[i].UserID == impression.UserID &&
			oldImpressions[i].AdID == impression.AdID {

			// Once we find the latest impression that has same userId and adId, check if it is within TTL
			if time.Since(oldImpressions[i].Timestamp) < cfg.TTL {
				return nil // Handling duplicate impression in service layer allows us to return custom error and then return custom response in handler
			}

			if time.Since(oldImpressions[i].Timestamp) >= cfg.TTL {
				break // If the latest impression is older than TTL, we can stop the loop and create a new impression
			}
		}
	}

	return storage.CreateImpression(impression)
}
