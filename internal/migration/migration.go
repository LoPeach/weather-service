package migration

import (
	"github.com/jmoiron/sqlx"
)

func InitSchema(db *sqlx.DB) {
	db.Exec(`CREATE TABLE if not exists weather (
    id SERIAL PRIMARY KEY,
    city TEXT,
    temperature DOUBLE PRECISION,
    created_at TIMESTAMP DEFAULT NOW()
);`)

	db.Exec(`ALTER TABLE weather 
ADD COLUMN IF NOT EXISTS 
measured_at TIMESTAMP;`)

	db.Exec(`ALTER TABLE weather ADD CONSTRAINT unique_city_measurement UNIQUE (city, measured_at);`)
}
