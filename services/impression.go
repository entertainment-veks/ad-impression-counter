package services

import (
	"ad_impression_counter/config"
	"ad_impression_counter/model"
	"ad_impression_counter/storage"
	"errors"
)

func TrackImpression(impression model.Impression, cfg config.Config) error {
	campaign, err := GetCampaign(impression.CampaignID)
	if err != nil {
		return errors.New("campaign not found")
	}

	if campaign.StartTime.After(impression.Timestamp) {
		return errors.New("campaign has not started yet")
	}

	return storage.SaveOrDiscardImpression(impression, cfg.TTL)
}
