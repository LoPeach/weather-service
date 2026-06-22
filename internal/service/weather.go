package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"new_begin/internal/client"
	"new_begin/internal/models"
	"new_begin/internal/repository"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrInvalidWeather = errors.New("invalid weather data")

var InvalidRequestToDB = errors.New("Ошибка запроса в БД")

var ErrWeatherAlreadyExists = errors.New(
	"weather already exists for this city in last 10 minutes",
)

func CreateWeather(ctx context.Context, db *sqlx.DB, city string, lat, lon float64) (models.Weather, error) {

	if city != "" {
		respWeather, err := client.GetWeatherByCity(ctx, city)
		if err != nil {
			return models.Weather{}, fmt.Errorf("get weather from api: %w", err)
		}
		if err = CheckAPIWeather(respWeather); err != nil {
			return models.Weather{}, fmt.Errorf("validate weather: %w", err)
		}

		if err = CheckTimeWeather(ctx, db, respWeather.City); err != nil {
			return models.Weather{}, fmt.Errorf("check cache time: %w", err)
		}

		if err = repository.SaveWeather(ctx, db, respWeather); err != nil {
			return models.Weather{}, fmt.Errorf("save weather: %w", err)
		}
		return respWeather, nil
	} else {
		respWeather, err := client.GetWeatherByCoordinates(ctx, lat, lon)
		if err != nil {
			return models.Weather{}, fmt.Errorf("get weather from api: %w", err)
		}
		if err = CheckAPIWeather(respWeather); err != nil {
			return models.Weather{}, fmt.Errorf("validate weather: %w", err)
		}

		if err = CheckTimeWeather(ctx, db, respWeather.City); err != nil {
			return models.Weather{}, fmt.Errorf("check cache time: %w", err)
		}

		if err = repository.SaveWeather(ctx, db, respWeather); err != nil {
			return models.Weather{}, fmt.Errorf("save weather: %w", err)
		}
		return respWeather, nil
	}
}

func CheckAPIWeather(inbox models.Weather) error {
	switch temp := inbox.Temperature; {
	case temp == 0.0 || temp == -999:
		return ErrInvalidWeather
	default:
		return nil
	}
}
func CheckTimeWeather(ctx context.Context, db *sqlx.DB, city string) error {
	allData, err := repository.GetLastWeatherByCity(ctx, db, city)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get last weather by city: %w", err)
	}
	firstDate := allData.CreatedAt
	if time.Since(firstDate) <= 10*time.Minute {
		return ErrWeatherAlreadyExists
	}
	return nil
}

func DeleteWeatherByCity(ctx context.Context, db *sqlx.DB, city models.WeatherCity) error {
	err := repository.DeleteWeatherByCity(ctx, db, city)
	if err != nil {
		log.Printf("Ошибка запроса %v", err)
		return err
	}
	return nil

}

func GetAllWeather(ctx context.Context, db *sqlx.DB) ([]models.Weather, error) {
	var w []models.Weather
	w, err := repository.GetAllWeather(db, ctx)
	if err != nil {
		log.Printf("Ошибка получения погоды %v", err)
		return w, err
	}
	return w, nil
}

func GetWeatherByCity(ctx context.Context, db *sqlx.DB, city string) ([]models.Weather, error) {
	var w []models.Weather
	w, err := repository.GetWeatherByCity(ctx, db, city)
	if err != nil {
		return w, err
	}
	return w, nil
}
