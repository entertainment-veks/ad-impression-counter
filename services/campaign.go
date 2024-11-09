package services

import (
	"ad_impression_counter/model"
	"ad_impression_counter/storage"
)

func CreateCampaign(campaign model.Campaign) error {
	return storage.CreateCampaign(campaign)
}

func GetCampaign(id string) (*model.Campaign, error) {
	campaign, err := storage.GetCampaignByID(id)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}
