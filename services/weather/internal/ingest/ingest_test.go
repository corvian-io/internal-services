package ingest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corvian/argus/services/weather/internal/config"
)

func TestFetch(t *testing.T) {
	resp := []byte(`{
    "request": {
        "type": "Zipcode",
        "query": "21117",
        "language": "en",
        "unit": "f"
    },
    "location": {
        "name": "Owings Mills",
        "country": "USA",
        "region": "Maryland",
        "lat": "39.427",
        "lon": "-76.777",
        "timezone_id": "America/New_York",
        "localtime": "2026-08-12 15:03",
        "localtime_epoch": 1786546980,
        "utc_offset": "-4.0"
    },
    "current": {
        "observation_time": "07:03 PM",
        "temperature": 90,
        "weather_code": 113,
        "weather_icons": [
            "https://cdn.worldweatheronline.com/images/wsymbols01_png_64/wsymbol_0001_sunny.png"
        ],
        "weather_descriptions": [
            "Sunny"
        ],
        "astro": {
            "sunrise": "06:17 AM",
            "sunset": "08:07 PM",
            "moonrise": "05:58 AM",
            "moonset": "08:14 PM",
            "moon_phase": "New Moon",
            "moon_illumination": 1
        },
        "air_quality": {
            "co": "120",
            "no2": "1.4",
            "o3": "121",
            "so2": "1.2",
            "pm2_5": "11.1",
            "pm10": "11.6",
            "us-epa-index": "1",
            "gb-defra-index": "1"
        },
        "wind_speed": 3,
        "wind_degree": 253,
        "wind_dir": "WSW",
        "pressure": 1011,
        "precip": 0,
        "humidity": 58,
        "cloudcover": 2,
        "feelslike": 97,
        "uv_index": 5,
        "visibility": 6,
        "is_day": "yes"
    }
}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}))
	defer server.Close()
	c := NewClient(config.Config{ApiKey: "test", Url: server.URL})
	weather, err := c.Fetch(context.Background())
	if err != nil {
		t.Fatalf("excepted weather; got %v", err)
	}
	t.Logf("got weather! %v", weather)
}
