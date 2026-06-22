package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"new_begin/internal/client"
	"new_begin/internal/handler"
	"new_begin/internal/middleware"
	"new_begin/internal/migration"
	"new_begin/internal/models"
	"new_begin/internal/repository"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	var finalResult []models.Weather
	ctx := context.Background()
	citiesSlice := []models.Cities{
		{Name: "Moscow", Latitude: 55.75, Longitude: 37.61},
		{Name: "Saratov", Latitude: 51.54, Longitude: 46.00},
		{Name: "Vladivostok", Latitude: 43.11, Longitude: 131.93},
	}

	for _, city := range citiesSlice {
		resp, err := client.GetWeatherByCity(ctx, city.Name)
		if err != nil {
			log.Printf("Ошибка запроса:%s", err)

			continue
		}
		finalResult = append(finalResult, resp)

	}
	for _, fn := range finalResult {
		fmt.Printf("В городе %s температура составляет %f градусов Цельсия\n", fn.City, fn.Temperature)
	}
	dsn := "postgres://weather_user:weather_pass@localhost:5432/weather_db?sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln("Не удалось подключиться к БД:", err)
	}
	defer db.Close()
	migration.InitSchema(db)

	for _, ins := range finalResult {
		err := repository.SaveWeather(ctx, db, ins)
		if err != nil {
			log.Printf("Ошибка выполнения вставки для города %s: %v", ins.City, err)
			continue
		}

	}
	fmt.Println("Connected successfully")
	var allWeather []models.Weather
	allWeather, err = repository.GetAllWeather(db, ctx)
	if err != nil {
		log.Printf("Ошибка получения всей погоды:%v", err)
	}
	fmt.Println(allWeather)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/weather", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.HandleGetWeather(db, w, r)
		case http.MethodPost:
			handler.HandlePostWeather(db, w, r)
		case http.MethodDelete:
			handler.HandleDeleteWeather(db, w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}

	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware.Logging(mux),
	}

	log.Println("Сервер запущен на http://localhost:8080/")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка сервера: %v", err)
	}

}
