package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"new_begin/internal/models"
	"time"
)

var key = "xxxxx"

func GetWeatherByCity(ctx context.Context, c string) (models.Weather, error) {

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?q=%s&key=%s", c, key)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return models.Weather{}, err
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.Weather{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return models.Weather{}, err
	}

	fmt.Println(string(body))

	var weather models.CurrentFromJSON

	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println("Ошибка парсинга Json:", err)
		return models.Weather{}, err
	}

	var finalWeather models.Weather

	finalWeather.City = c
	finalWeather.Temperature = weather.WeatherFromJSON.WeatherTemperature
	loc, err := time.LoadLocation(
		weather.LocationFromJSON.TZID,
	)
	if err != nil {
		return models.Weather{}, err
	}

	finalWeather.MeasuredAt, err = time.ParseInLocation(
		"2006-01-02 15:04",
		weather.WeatherFromJSON.MeasuredAt,
		loc,
	)

	if err != nil {
		return models.Weather{}, err
	}

	finalWeather.MeasuredAt = finalWeather.MeasuredAt.UTC()

	return finalWeather, nil

}
func GetWeatherByCoordinates(ctx context.Context, lat float64, lon float64) (models.Weather, error) {
	url := fmt.Sprintf(
		"https://api.weatherapi.com/v1/current.json?q=%.2f,%.2f&key=%s",
		lat,
		lon,
		key,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return models.Weather{}, err
	}

	tr := &http.Transport{
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	resp, err := client.Do(req)

	if err != nil {
		return models.Weather{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return models.Weather{}, err
	}

	var weather models.CurrentFromJSON

	err = json.Unmarshal(body, &weather)
	if err != nil {
		return models.Weather{}, err
	}

	var finalWeather models.Weather

	finalWeather.City = weather.LocationFromJSON.City
	finalWeather.Temperature = weather.WeatherFromJSON.WeatherTemperature

	loc, err := time.LoadLocation(
		weather.LocationFromJSON.TZID,
	)
	if err != nil {
		return models.Weather{}, err
	}

	finalWeather.MeasuredAt, err = time.ParseInLocation(
		"2006-01-02 15:04",
		weather.WeatherFromJSON.MeasuredAt,
		loc,
	)
	if err != nil {
		return models.Weather{}, err
	}

	finalWeather.MeasuredAt = finalWeather.MeasuredAt.UTC()

	return finalWeather, nil
}
