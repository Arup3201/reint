package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"
)

const (
	host = "localhost"
	port = "8080"
)

type WindData struct {
	TargetTime  string   `json:"time"`
	PublishTime *string  `json:"publish_time"`
	Actual      float64  `json:"actual"`
	Forecast    *float64 `json:"forecast"`
}

var (
	DATASET = []WindData{}
)

func LoadDataset() error {

	windApi := "https://data.elexon.co.uk/bmrs/api/v1/datasets/FUELHH/stream?settlementDateFrom=2024-01-01&settlementDateTo=2024-01-31&fuelType=WIND"
	forecastApi := "https://data.elexon.co.uk/bmrs/api/v1/datasets/WINDFOR/stream?publishDateTimeFrom=2024-01-01&publishDateTimeTo=2024-01-31"

	type apiData struct {
		StartTime   string  `json:"startTime"`
		PublishTime string  `json:"publishTime"`
		Generation  float64 `json:"generation"` // MW
	}

	var req *http.Request
	var res *http.Response
	var err error

	req, _ = http.NewRequest("GET", windApi, nil)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("actual wind data response: %w", err)
	}
	defer res.Body.Close()

	actualWindData := []apiData{}
	if err = json.NewDecoder(res.Body).Decode(&actualWindData); err != nil {
		return fmt.Errorf("actual wind data response decode: %w", err)
	}

	req, _ = http.NewRequest("GET", forecastApi, nil)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("forecasted wind data response: %w", err)
	}
	defer res.Body.Close()

	forecastedWindData := []apiData{}
	if err := json.NewDecoder(res.Body).Decode(&forecastedWindData); err != nil {
		return fmt.Errorf("forecasted wind data response decode: %w", err)
	}

	DATASET = []WindData{}
	var isNull bool
	for _, actual := range actualWindData {
		isNull = true
		for _, forecast := range forecastedWindData {
			if actual.StartTime == forecast.StartTime {
				isNull = false
				DATASET = append(DATASET, WindData{
					TargetTime:  actual.StartTime,
					PublishTime: &forecast.PublishTime,
					Actual:      actual.Generation,
					Forecast:    &forecast.Generation,
				})
			}
		}

		if isNull {
			DATASET = append(DATASET, WindData{
				TargetTime: actual.StartTime,
				Actual:     actual.Generation,
			})
		}
	}

	var parseTimeA, parseTimeB time.Time
	slices.SortFunc(DATASET, func(a, b WindData) int {
		parseTimeA, _ = time.Parse(time.RFC3339, a.TargetTime)
		parseTimeB, _ = time.Parse(time.RFC3339, b.TargetTime)
		return int(parseTimeA.Sub(parseTimeB))
	})

	return nil
}

func getWindData(w http.ResponseWriter, r *http.Request) {
	var err error

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	horizon := r.URL.Query().Get("horizon")

	if start == "" || end == "" || horizon == "" {
		http.Error(w, "Start, end or horizon value cannot be empty", http.StatusBadRequest)
		return
	}

	var minTime, maxTime time.Time
	minTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	maxTime = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	var startTime time.Time
	startTime, err = time.Parse(time.RFC3339, start)
	if err != nil {
		log.Printf("[ERROR] parse start time: %s\n", err)
		http.Error(w, "Start time is not valid, valid format 2024-01-01T02:30:00Z", http.StatusBadRequest)
		return
	}

	var endTime time.Time
	endTime, err = time.Parse(time.RFC3339, end)
	if err != nil {
		log.Printf("[ERROR] parse end time: %s\n", err)
		http.Error(w, "End time is not valid, valid format 2024-01-01T02:30:00Z", http.StatusBadRequest)
		return
	}

	if minTime.Sub(startTime) > 0 || startTime.Sub(maxTime) >= 0 ||
		minTime.Sub(endTime) > 0 || endTime.Sub(maxTime) >= 0 {
		http.Error(w, "Only January 2024 values are accepted for start time and end time", http.StatusBadRequest)
		return
	}

	if endTime.Sub(startTime) < 0 {
		http.Error(w, "End time cannot be less than start time", http.StatusBadRequest)
		return
	}

	var horizonInt int
	horizonInt, err = strconv.Atoi(horizon)
	if err != nil {
		log.Printf("[ERROR] parse horizon: %s\n", err)
		http.Error(w, "Forecast horizon should be an integer in range 0-48 Hrs", http.StatusBadRequest)
		return
	}

	if horizonInt < 0 || horizonInt > 48 {
		http.Error(w, "Forecast horizon should be an integer in range 0-48 Hrs", http.StatusBadRequest)
		return
	}

	var parsedTargetTime, parsedPublishTime time.Time
	var results = []WindData{}
	for _, data := range DATASET {
		parsedTargetTime, _ = time.Parse(time.RFC3339, data.TargetTime)
		if data.PublishTime != nil {
			parsedPublishTime, _ = time.Parse(time.RFC3339, *data.PublishTime)
		}
		if parsedTargetTime.Sub(startTime) >= 0 &&
			endTime.Sub(parsedTargetTime) >= 0 &&
			(data.PublishTime == nil ||
				parsedTargetTime.Sub(parsedPublishTime) >= time.Duration(horizonInt)*time.Hour) {
			results = append(results, data)
		}
	}

	if len(results) == 0 {
		json.NewEncoder(w).Encode(map[string]any{
			"data": results,
		})
		return
	}

	var filteredResults []WindData
	i, k := 0, 0

	filteredResults = append(filteredResults, results[i])
	i += 1
	k += 1
	for ; i < len(results); i++ {
		if filteredResults[k-1].TargetTime == results[i].TargetTime {
			filteredResults[k-1] = results[i]
		} else {
			filteredResults = append(filteredResults, results[i])
			k += 1
		}
	}

	json.NewEncoder(w).Encode(map[string]any{
		"data": filteredResults,
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	var err error

	err = LoadDataset()
	if err != nil {
		log.Fatalf("load dataset: %s\n", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /wind-data", getWindData)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      CorsMiddleware(mux),
	}

	log.Printf("[INFO] Server starting at %s:%s\n", host, port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("server listen and serve: %s\n", err)
	}
}
