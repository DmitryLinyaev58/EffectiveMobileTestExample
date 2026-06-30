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

// handleGetSubscriptionByID — получить подписку по ID
// @Summary Получить подписку по ID
// @Description Получить подробную информацию о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} models.Subscription
// @Failure 404 {string} string "Подписка не найдена"
// @Router /subscriptions/{id} [get]
func handleGetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	// Оставляем твой ручной парсинг пути — он рабочий
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
}

// handleCreateSubscription — создать подписку
// @Summary Создать подписку
// @Description Создать новую подписку в системе
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param body body db.CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} db.Subscription
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /subscriptions [post]
func handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req db.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Простая валидация
	if req.ServiceName == "" {
		http.Error(w, "service_name is required", http.StatusBadRequest)
		return
	}

	repo := db.NewSubscriptionRepository(database)

	sub, err := repo.Create(r.Context(), req.UserID, req.ServiceName, req.Price, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("❌ DB error on create: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

func Run() error {
	cfg := config.Load()
	log.Printf("📋 Конфигурация: host=%s, port=%s, db=%s", cfg.DBHost, cfg.Port, cfg.DBName)

	var err error
	database, err = db.Connect(cfg)
	if err != nil {
		return fmt.Errorf("❌ не удалось подключиться к БД: %w", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("EffectiveMobile API is running"))
	})

	// --- Swagger UI (HTML) ---
	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>Swagger UI — Effective Mobile</title>
    <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui-bundle.js"></script>
    <script>
      window.onload = function() {
        const ui = SwaggerUIBundle({
          url: "/swagger/spec",
          dom_id: '#swagger-ui',
          presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIBundle.SwaggerUIStandalonePreset
          ],
          layout: "BaseLayout"
        })
        window.ui = ui
      }
    </script>
  </body>
</html>`
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	})

	// --- Swagger JSON (спецификация) ---
	mux.HandleFunc("/swagger/spec", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	mux.HandleFunc("/subscriptions/", handleGetSubscriptionByID)

	mux.HandleFunc("/subscriptions", handleCreateSubscription)

	addr := ":" + cfg.Port
	log.Printf("🚀 Сервер запущен на порту :%s", addr)

	return http.ListenAndServe(addr, mux)
}
