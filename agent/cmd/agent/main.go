package main

import (
    "log"
    "time"

    "Kvision/agent/pkg/config"
    "Kvision/agent/pkg/collector"
    "Kvision/agent/pkg/sender"
)


func main() {
	log.Println("🚀 Starting KubeVision Agent...")

	// 1. Load configuration (PROM_URL, BACKEND_URL, SCRAPE_INTERVAL)
	cfg := config.Load()
	log.Printf("Loaded config: PromURL=%s BackendURL=%s Interval=%s",
		cfg.PromURL, cfg.BackendURL, cfg.ScrapeInterval)

	// 2. Create a ticker that runs every SCRAPE_INTERVAL
	ticker := time.NewTicker(cfg.ScrapeInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		// 3. Scrape Prometheus
		metrics, err := collector.ScrapePrometheus(cfg.PromURL)
		if err != nil {
			log.Printf("❌ Error scraping Prometheus: %v", err)
			continue
		}

		log.Printf("Scraped %d metric samples", len(metrics))

		// 4. Push metrics to backend
		err = sender.PushMetrics(cfg.BackendURL, metrics)
		if err != nil {
			log.Printf("❌ Error sending metrics to backend: %v", err)
			continue
		}

		log.Println("✅ Successfully pushed metrics to backend")
	}
}
