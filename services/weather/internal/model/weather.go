package model

type location struct {
	Name    string
	Country string
	Region  string
}

type Weather struct {
	Location    location
	Temp        int
	Description string
	Icon        string
}
