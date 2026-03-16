package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestLoadDataset(t *testing.T) {
// 	err := LoadDataset()

// 	assert.NoError(t, err)
// 	assert.NotEqual(t, 0, len(DATASET))

// 	var cnt int

// 	cnt = 0
// 	for _, data := range DATASET {
// 		if data.TargetTime == "2024-01-03T21:00:00Z" {
// 			cnt += 1
// 		}
// 	}
// 	assert.Equal(t, 24, cnt)

// 	cnt = 0
// 	for _, data := range DATASET {
// 		if data.TargetTime == "2024-01-30T21:00:00Z" {
// 			cnt += 1
// 		}
// 	}
// 	assert.Equal(t, 24, cnt)
// }

func TestWindData_GET(t *testing.T) {
	log.SetOutput(io.Discard)

	DATASET = []WindData{
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-03T21:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
	}

	t.Run("should get wind data", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 200, res.StatusCode)
	})

	DATASET = []WindData{
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T21:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T21:30:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T22:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T22:30:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T23:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
		{
			PublishTime: &[]string{"2024-01-01T03:00:00Z"}[0],
			TargetTime:  "2024-01-01T23:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
	}

	t.Run("should get 5 results between start on 1st Jan and end on 2nd Jan with horizon 4", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		var responseBody struct {
			Data []WindData `json:"data"`
		}
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			assert.NoError(t, err)
			return
		}
		assert.Equal(t, 5, len(responseBody.Data))
	})

	t.Run("should get 2 results between start on 1st Jan 21:00 to 21:30 with horizon 4", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T21:00:00Z&end=2024-01-01T21:30:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		var responseBody struct {
			Data []WindData `json:"data"`
		}
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			assert.NoError(t, err)
			return
		}
		assert.Equal(t, 2, len(responseBody.Data))
	})

	DATASET = []WindData{
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T12:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
		{
			PublishTime: &[]string{"2024-01-01T03:00:00Z"}[0],
			TargetTime:  "2024-01-01T12:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
	}
	t.Run("should get 1 result with latest forecast at same target time", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		var responseBody struct {
			Data []WindData `json:"data"`
		}
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			assert.NoError(t, err)
			return
		}
		assert.Equal(t, 1, len(responseBody.Data))
	})

	DATASET = []WindData{
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T12:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
		{
			PublishTime: &[]string{"2024-01-01T03:00:00Z"}[0],
			TargetTime:  "2024-01-01T12:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T12:30:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
	}

	t.Run("should get 2 results with latest forecast with target published at 2 different times", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		var responseBody struct {
			Data []WindData `json:"data"`
		}
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			assert.NoError(t, err)
			return
		}
		assert.Equal(t, 2, len(responseBody.Data))
	})

	DATASET = []WindData{
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T12:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
		{
			PublishTime: &[]string{"2024-01-01T02:00:00Z"}[0],
			TargetTime:  "2024-01-01T12:30:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T13:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
	}
	t.Run("should get 2 results with horizon 10", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z&horizon=10", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		var responseBody struct {
			Data []WindData `json:"data"`
		}
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			assert.NoError(t, err)
			return
		}
		assert.Equal(t, 2, len(responseBody.Data))
	})

	DATASET = []WindData{
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T12:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
		{
			TargetTime: "2024-01-01T12:30:00Z",
			Actual:     6945,
		},
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-01T13:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
	}

	t.Run("should get 3 results", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		var responseBody struct {
			Data []WindData `json:"data"`
		}
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			assert.NoError(t, err)
			return
		}
		assert.Equal(t, 3, len(responseBody.Data))
	})

	DATASET = []WindData{
		{
			PublishTime: &[]string{"2024-01-01T02:30:00Z"}[0],
			TargetTime:  "2024-01-03T21:00:00Z",
			Actual:      6945,
			Forecast:    &[]float64{10629}[0],
		},
	}

	t.Run("should get bad request without start time", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?end=2024-01-02T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("should get bad request without end time", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("should get bad request without forecast horizon", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("should get bad request for invalid start time", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00Z&end=2024-01-02T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("should get bad request for start time in Feb 2024", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-02-01T00:00Z&end=2024-01-02T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("should get bad request for start time in Dec 2023", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2023-12-30T00:00Z&end=2024-01-02T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("should get bad request for invalid end time", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2024-01-02T00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("should get bad request for end time in Feb 2024", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2024-02-02T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("should get bad request for end time in Dec 2023", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2023-12-30T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("should get bad request for invalid forecast horizon", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z&horizon=y", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("should get bad request for forecast horizon 60", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z&horizon=60", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("should get bad request for start time more than end time", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wind-data?start=2024-01-03T00:00:00Z&end=2024-01-02T00:00:00Z&horizon=4", nil)
		w := httptest.NewRecorder()
		getWindData(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 400, res.StatusCode)
	})
}
