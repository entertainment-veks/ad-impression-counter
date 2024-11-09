package service

import (
	"ad_impression_counter/model"
	"ad_impression_counter/storage"
	"errors"
	"log"
)

var (
	ErrCampaignNotStarted = errors.New("campaign has not started yet")
)

func CreateCampaign(campaign model.Campaign) error {
	return storage.CreateCampaign(campaign)
}

func GetCampaign(id string) (*model.Campaign, error) {
	campaign, err := storage.GetCampaignByID(id)
	if err != nil {
		log.Printf("failed to get campaign by ID: %s: %v", id, err)
		return nil, err
	}

	return &campaign, nil
}
