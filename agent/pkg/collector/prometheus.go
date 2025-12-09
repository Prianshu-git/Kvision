package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)


// Metric Data Structure
type MetricSample struct {
	Namespace string  `json:"namespace"`
	Pod       string  `json:"pod"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Restarts  int     `json:"restarts"`
	Timestamp int64   `json:"timestamp"`
}

// Prometheus Response Types
type PromResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// Prometheus Query Function
func queryPrometheus(promURL string, promQL string) (*PromResponse, error) {
	fullURL := fmt.Sprintf("%s/api/v1/query?query=%s", promURL, promQL)

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error calling Prometheus API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response from Prometheus: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body) // updated here
	if err != nil {
		return nil, fmt.Errorf("error reading Prometheus response: %w", err)
	}

	var promResp PromResponse
	err = json.Unmarshal(body, &promResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing Prometheus JSON: %w", err)
	}

	return &promResp, nil
}

// Main Scraper
func ScrapePrometheus(promURL string) ([]MetricSample, error) {
	timestamp := time.Now().Unix()

	// PromQL queries
	cpuQuery := `sum(rate(container_cpu_usage_seconds_total{container!="",namespace!~"kube-system|monitoring"}[5m])) by (namespace, pod)`
	memQuery := `sum(container_memory_working_set_bytes{container!="",namespace!~"kube-system|monitoring"}) by (namespace, pod)`
	restartQuery := `sum(kube_pod_container_status_restarts_total{namespace!~"kube-system|monitoring"}) by (namespace, pod)`

	// Query Prometheus
	cpuResp, err := queryPrometheus(promURL, cpuQuery)
	if err != nil {
		log.Printf("Error querying CPU: %v", err)
	}

	memResp, err := queryPrometheus(promURL, memQuery)
	if err != nil {
		log.Printf("Error querying Memory: %v", err)
	}

	restartResp, err := queryPrometheus(promURL, restartQuery)
	if err != nil {
		log.Printf("Error querying Restarts: %v", err)
	}

	metricsMap := map[string]*MetricSample{}

	// Helper to generate key
	key := func(ns, pod string) string { return ns + "/" + pod }

	// CPU
	for _, r := range cpuResp.Data.Result {
		ns := r.Metric["namespace"]
		pod := r.Metric["pod"]
		cpuVal := parsePromValue(r.Value)

		k := key(ns, pod)
		if metricsMap[k] == nil {
			metricsMap[k] = &MetricSample{
				Namespace: ns,
				Pod:       pod,
				Timestamp: timestamp,
			}
		}
		metricsMap[k].CPU = cpuVal
	}

	// Memory
	for _, r := range memResp.Data.Result {
		ns := r.Metric["namespace"]
		pod := r.Metric["pod"]
		memVal := parsePromValue(r.Value)

		k := key(ns, pod)
		if metricsMap[k] == nil {
			metricsMap[k] = &MetricSample{
				Namespace: ns,
				Pod:       pod,
				Timestamp: timestamp,
			}
		}
		metricsMap[k].Memory = memVal
	}

	// Restarts
	for _, r := range restartResp.Data.Result {
		ns := r.Metric["namespace"]
		pod := r.Metric["pod"]
		restartVal := int(parsePromValue(r.Value))

		k := key(ns, pod)
		if metricsMap[k] == nil {
			metricsMap[k] = &MetricSample{
				Namespace: ns,
				Pod:       pod,
				Timestamp: timestamp,
			}
		}
		metricsMap[k].Restarts = restartVal
	}

	// Convert map → slice
	var result []MetricSample
	for _, m := range metricsMap {
		result = append(result, *m)
	}

	return result, nil
}

// Helper to parse Prometheus values
func parsePromValue(v []interface{}) float64 {
	if len(v) != 2 {
		return 0
	}
	strVal, ok := v[1].(string)
	if !ok {
		return 0
	}
	var parsed float64
	fmt.Sscan(strVal, &parsed)
	return parsed
}
