package model

type Location struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Region  string `json:"region"`
}

type CurrentWeather struct {
	Temp        int      `json:"temperature"`
	Description []string `json:"weather_descriptions"`
	Icon        []string `json:"weather_icons"`
}

type Weather struct {
	Location Location       `json:"location"`
	Current  CurrentWeather `json:"current"`
}
