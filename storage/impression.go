package storage

import (
	"ad_impression_counter/model"
	"errors"
	"sync"
)

var (
	// It is a map that stores impressions by campaign ID
	// Value is a slice of impressions, sorted by timestamp
	impressions = sync.Map{} //map[campaignID][]model.Impression
)

func CreateImpression(impression model.Impression) error {
	rawImpressions, exists := impressions.Load(impression.CampaignID)
	if !exists {
		impressions.Store(impression.CampaignID, []model.Impression{impression})
		return nil
	}

	oldImpressions, ok := rawImpressions.([]model.Impression)
	if !ok {
		return errors.New("invalid old impression data")
	}

	impressions.Store(impression.CampaignID, append(oldImpressions, impression))
	return nil
}

func GetImpressionsByCampaign(campaignID string) ([]model.Impression, error) {
	rawImpressions, exists := impressions.Load(campaignID)
	if !exists {
		return nil, nil
	}

	impressions, ok := rawImpressions.([]model.Impression)
	if !ok {
		return nil, errors.New("invalid impression data")
	}

	return impressions, nil
}
