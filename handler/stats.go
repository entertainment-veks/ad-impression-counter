package handler

import (
	"ad_impression_counter/services"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterStatsRoutes(router *mux.Router) {
	router.HandleFunc("/campaigns/{id}/stats", getCampaignStatsHandler).Methods("GET")
}

func getCampaignStatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignID := vars["id"]

	if campaignID == "" {
		http.Error(w, "Campaign ID is required", http.StatusBadRequest)
		return
	}

	stats, err := services.GetCampaignStats(campaignID)
	if errors.Is(err, services.ErrCampaignNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
