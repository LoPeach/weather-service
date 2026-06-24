package main

import (
	"fmt"
	"log"
	"net/http"
	"new_begin/internal/handler"
	"new_begin/internal/middleware"
	"new_begin/internal/migration"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func waitForDB(dsn string) *sqlx.DB {
	for {
		db, err := sqlx.Open("pgx", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("✅ database is ready")
				return db
			}
			db.Close()
		}

		log.Println("⏳ waiting for database...")
		time.Sleep(2 * time.Second)
	}
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	log.Println("DSN =", dsn)

	db := waitForDB(dsn)
	defer db.Close()

	migration.InitSchema(db)

	fmt.Println("Connected successfully")

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
