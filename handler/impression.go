package handler

import (
	"ad_impression_counter/config"
	"ad_impression_counter/model"
	"ad_impression_counter/service"
	"ad_impression_counter/storage"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type TrackImpressionRequest struct {
	CampaignID string `json:"campaign_id"`
	UserID     string `json:"user_id"`
	AdID       string `json:"ad_id"`
}

func RegisterImpressionRoutes(router *mux.Router, cfg config.Config) {
	router.HandleFunc("/impressions", trackImpressionHandler(cfg)).Methods("POST")
}

func trackImpressionHandler(cfg config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TrackImpressionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.CampaignID == "" {
			http.Error(w, "Campaign ID is required", http.StatusBadRequest)
			return
		}

		if req.UserID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		if req.AdID == "" {
			http.Error(w, "Ad ID is required", http.StatusBadRequest)
			return
		}

		impression := model.Impression{
			CampaignID: req.CampaignID,
			UserID:     req.UserID,
			AdID:       req.AdID,
			Timestamp:  time.Now(),
		}

		err := service.TrackImpression(impression, cfg)
		if errors.Is(err, storage.ErrCampaignNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrCampaignNotStarted) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Impression tracked successfully"})
	}
}
