package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"new_begin/internal/models"
	"new_begin/internal/response"
	"new_begin/internal/service"

	"github.com/jmoiron/sqlx"
)

func HandleGetWeather(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	city := queryParams.Get("city")
	ctx := r.Context()
	data, err := service.GetAllWeather(ctx, db)
	if err != nil {
		log.Printf("Ошибка получения погоды %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	var filterData []models.Weather
	if len(city) > 0 {
		filterData, err = service.GetWeatherByCity(ctx, db, city)
		if err != nil {
			log.Printf("Ошибка получения погоды %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

	}

	if city != "" {
		response.Success(w, http.StatusOK, filterData)
		return
	}
	response.Success(w, http.StatusOK, data)

}

func HandlePostWeather(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	var req models.CreateWeatherRequest
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	if req.City == "" && req.Latitude == nil && req.Longitude == nil {
		response.Error(w, http.StatusBadRequest, "Название города или пара долгота/широта должны быть заполнены")
		return
	}

	if req.Latitude == nil && req.Longitude != nil {
		response.Error(w, http.StatusBadRequest, "Широта для пары долгота/широта обязательна для заполнения")
		return
	}

	if req.Longitude == nil && req.Latitude != nil {
		response.Error(w, http.StatusBadRequest, "Долгота для пары долгота/широта обязательна для заполнения")
		return
	}

	var lat float64
	var lng float64

	if req.Latitude != nil {
		lat = *req.Latitude
	}

	if req.Longitude != nil {
		lng = *req.Longitude
	}

	_, err = service.CreateWeather(ctx, db, req.City, lat, lng)
	if errors.Is(err, service.ErrInvalidWeather) {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	} else if errors.Is(err, service.ErrWeatherAlreadyExists) {
		response.Error(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		log.Printf("POST ERROR: %+v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.Success(w, http.StatusCreated, "weather created successfully")

}

func HandleDeleteWeather(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	var req models.WeatherCity
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.City == "" {
		response.Error(w, http.StatusBadRequest, "Имя обязательно для заполнения")
		return
	}

	err = service.DeleteWeatherByCity(ctx, db, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Ошибка удаление погоды для города")
		return
	}

	response.Success(w, http.StatusOK, "Weather deleted successfully")

}
