package storage

import (
	"ad_impression_counter/model"
	"errors"
	"sync"
	"time"
)

var (
	impressions = sync.Map{}
)

func SaveOrDiscardImpression(newImpression model.Impression, ttl time.Duration) error {
	key := GetKey(newImpression.CampaignID, newImpression.UserID, newImpression.AdID)

	rawImpression, exists := impressions.Load(key)
	if exists { // If the impression already exists, check timestamp
		oldImpression, ok := rawImpression.(model.Impression)
		if !ok {
			return errors.New("invalid old impression data")
		}

		if time.Since(oldImpression.Timestamp) < ttl { // If timestamp is within TTL, discard
			return nil
		}
	}

	impressions.Store(key, newImpression)
	return nil
}

func GetImpressionsByCampaign(campaignID string) []model.Impression {
	var impressionsByCampaign []model.Impression

	impressions.Range(func(key, value interface{}) bool {
		impression, ok := value.(model.Impression)
		if !ok {
			return true
		}

		if impression.CampaignID == campaignID {
			impressionsByCampaign = append(impressionsByCampaign, impression)
		}

		return true
	})

	return impressionsByCampaign
}

func GetKey(campaignId string, userId string, adId string) string {
	return campaignId + userId + adId
}
