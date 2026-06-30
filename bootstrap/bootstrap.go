package bootstrap

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"encoding/json"
	"strconv"
	"strings"

	"github.com/DmitryLinyaev58/EffectiveMobileTestExample/config"
	"github.com/DmitryLinyaev58/EffectiveMobileTestExample/db"
)

var database *sql.DB

func Run() error {
	cfg := config.Load()
	log.Printf("конфигурация: host=%s, port=%s, db=%s", cfg.DBHost, cfg.Port, cfg.DBName)

	var err error
	database, err = db.Connect(cfg)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("EffectiveMobile API is running"))
	})

	// =========================================================
	// ручки GET
	// =========================================================
	// Хендлер для GET /subscriptions/<id>
	mux.HandleFunc("/subscriptions/", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что после /subscriptions/ что-то есть
		path := strings.TrimPrefix(r.URL.Path, "/subscriptions/")
		if path == "" || strings.Contains(path, "/") {
			http.Error(w, "Invalid URL format. Expected /subscriptions/{id}", http.StatusBadRequest)
			return
		}

		idStr := path
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("❌ Invalid ID format: %v", err)
			http.Error(w, "ID must be a valid integer", http.StatusBadRequest)
			return
		}

		repo := db.NewSubscriptionRepository(database)
		sub, err := repo.GetByID(r.Context(), id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "Subscription not found", http.StatusNotFound)
			} else {
				log.Printf("❌ Database error: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sub)
	})
	// =========================================================
	//ручки END
	// =========================================================

	addr := ":" + cfg.Port
	log.Printf("🚀 Сервер запущен на порту :%s", addr)

	return http.ListenAndServe(addr, mux)
}
