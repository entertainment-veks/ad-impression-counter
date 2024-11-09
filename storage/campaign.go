package storage

import (
	"ad_impression_counter/model"
	"errors"
	"sync"
)

var (
	ErrCampaignNotFound = errors.New("campaign not found")

	campaigns = sync.Map{}
)

func CreateCampaign(campaign model.Campaign) error {
	campaigns.Store(campaign.ID, campaign)
	return nil
}

func GetCampaignByID(id string) (model.Campaign, error) {
	rawCampaign, exists := campaigns.Load(id)
	if !exists {
		return model.Campaign{}, ErrCampaignNotFound
	}

	campaign, ok := rawCampaign.(model.Campaign)
	if !ok {
		return model.Campaign{}, errors.New("invalid campaign data")
	}

	return campaign, nil
}
