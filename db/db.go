package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/DmitryLinyaev58/EffectiveMobileTestExample/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(cfg config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ошибка проверки связи (Ping): %w", err)
	}

	log.Println("✅ База данных подключена успешно!")
	return db, nil
}
