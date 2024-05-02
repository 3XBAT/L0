package repository

import (

	"fmt"

	"github.com/jmoiron/sqlx"

)

const (
	deliveryTable = "delivery"
	itemsTable = "items"
	paymentTable = "payment"
	ordersTable = "orders"
)


type Config struct {
	Port     string
	Username string
	Host     string
	DBName   string
	Password string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)) //
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {//проверяем можем ли подключться к нащей бд
		return nil, err
	}

	return db, nil

}
