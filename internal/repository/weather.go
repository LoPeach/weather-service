package repository

import (
	"context"
	"fmt"
	"new_begin/internal/models"

	"github.com/jmoiron/sqlx"
)

func GetAllWeather(db *sqlx.DB, ctx context.Context) ([]models.Weather, error) {
	var w []models.Weather
	query := `SELECT DISTINCT ON (city)
       city,
       temperature,
       measured_at,
       created_at
FROM weather
ORDER BY city, created_at DESC`

	err := db.Select(&w, query)
	if err != nil {
		return w, err
	}
	return w, nil

}

func GetWeatherByCity(ctx context.Context, db *sqlx.DB, city string) ([]models.Weather, error) {
	var w []models.Weather
	query := `
SELECT city, temperature, measured_at, created_at
FROM weather
WHERE city = $1
ORDER BY created_at DESC
LIMIT 1
`

	err := db.SelectContext(ctx, &w, query, city)
	if err != nil {
		return w, err
	}
	return w, nil

}

func SaveWeather(ctx context.Context, db *sqlx.DB, wthr models.Weather) error {
	fmt.Printf("%+v\n", wthr)
	_, err := db.NamedExecContext(ctx, `INSERT INTO weather (city, temperature, measured_at) 
VALUES (:city, :temperature, :measured_at)
ON CONFLICT (city, measured_at) 
DO UPDATE SET 
temperature = EXCLUDED.temperature,
created_at = NOW(); `, wthr)
	if err != nil {
		return err
	}
	return nil
}

func DeleteWeatherByCity(ctx context.Context, db *sqlx.DB, wthr models.WeatherCity) error {
	_, err := db.NamedExecContext(ctx, `DELETE FROM weather WHERE city = :city`, wthr)
	if err != nil {
		return err
	}
	return nil
}

func GetLastWeatherByCity(ctx context.Context, db *sqlx.DB, city string) (models.Weather, error) {
	var w models.Weather
	query := "SELECT city, temperature, measured_at, created_at FROM weather WHERE city =$1 ORDER BY created_at DESC LIMIT 1"

	err := db.GetContext(ctx, &w, query, city)
	if err != nil {
		return w, err
	}
	return w, nil

}
