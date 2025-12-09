package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"Kvision/agent/pkg/collector"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func PushMetrics(backendURL string, metrics []collector.MetricSample) error {
	if len(metrics) == 0 {
		log.Println("⚠️ No metrics to send, skipping push")
		return nil
	}

	// Convert metrics slice → JSON
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	// Final URL
	url := fmt.Sprintf("%s/ingest/metrics", backendURL)

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}
	defer resp.Body.Close()

	// Non-200 responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("backend returned non-200: %d", resp.StatusCode)
	}

	log.Printf("📤 Pushed %d metrics to backend", len(metrics))
	return nil
}
