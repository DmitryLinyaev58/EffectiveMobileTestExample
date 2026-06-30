package bootstrap

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DmitryLinyaev58/EffectiveMobileTestExample/config"
	"github.com/DmitryLinyaev58/EffectiveMobileTestExample/db"
)

func Run() error {

	cfg := config.Load()
	log.Printf("...нфигурация: host=%s, port=%s, db=%s", cfg.DBHost, cfg.Port, cfg.DBName)

	// 2. Подключение к БД
	database, err := db.Connect(cfg)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}
	defer database.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("EffectiveMobile API is running"))
	})

	addr := ":" + cfg.Port
	log.Printf("🚀 Сервер запущен на порту :%s", addr)

	return http.ListenAndServe(addr, mux)
}
