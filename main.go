package main

import (
	"ad_impression_counter/config"
	"ad_impression_counter/handler"
	"ad_impression_counter/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %s\n", err.Error())
	}

	startImpressionWorkers(cfg)

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	apiV1 := api.PathPrefix("/v1").Subrouter()

	handler.RegisterMiddlewares(apiV1)

	handler.RegisterCampaignRoutes(apiV1)
	handler.RegisterImpressionRoutes(apiV1, cfg)
	handler.RegisterStatsRoutes(apiV1)

	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

func startImpressionWorkers(cfg config.Config) {
	for i := 0; i < cfg.ImpressionWorkers; i++ {
		go func() {
			service.ConsumeAndProcessImpressions(cfg)
		}()
	}
}
