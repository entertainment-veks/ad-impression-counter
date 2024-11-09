package handler

import (
	"ad_impression_counter/model"
	"ad_impression_counter/service"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CreateCampaignRequest struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
}

func RegisterCampaignRoutes(router *mux.Router) {
	router.HandleFunc("/campaigns", createCampaignHandler).Methods("POST")
}

func createCampaignHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if req.StartTime.IsZero() {
		http.Error(w, "Start time is required", http.StatusBadRequest)
		return
	}

	if req.StartTime.Before(time.Now()) {
		http.Error(w, "Start time must be in the future", http.StatusBadRequest)
		return
	}

	campaign := model.Campaign{
		ID:        uuid.New().String(),
		Name:      req.Name,
		StartTime: req.StartTime,
	}

	if err := service.CreateCampaign(campaign); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": campaign.ID, "message": "Campaign created successfully"})
}
