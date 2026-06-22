package models

import "time"

type CreateWeatherRequest struct {
	City      string   `json:"city"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type ResponseInfo struct {
	Message string `json:"message"`
}

type WeatherCity struct {
	City string `db:"city"`
}
type Weather struct {
	City        string    `db:"city"`
	Temperature float64   `db:"temperature"`
	MeasuredAt  time.Time `db:"measured_at"`
	CreatedAt   time.Time `db:"created_at"`
}

type Cities struct {
	Name      string
	Latitude  float64
	Longitude float64
}

type CurrentFromJSON struct {
	WeatherFromJSON struct {
		MeasuredAt         string  `json:"last_updated"`
		WeatherTemperature float64 `json:"temp_c"`
	} `json:"current"`
	LocationFromJSON struct {
		City string `json:"name"`
		TZID string `json:"tz_id"`
	} `json:"location"`
}
