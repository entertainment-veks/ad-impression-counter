package service

import (
	"ad_impression_counter/config"
	"ad_impression_counter/model"
	"ad_impression_counter/storage"
	"log"
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

	return storage.SaveOrDiscardImpression(impression, cfg.TTL)
}
