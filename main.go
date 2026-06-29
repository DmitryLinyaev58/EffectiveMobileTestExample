package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
)

// Config храним все настройки приложения
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
}

// loadConfig собирает конфигурацию из переменных окружения
// Если переменная не найдена, подставляется дефолтное значение
func loadConfig() Config {
	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "effective_mobile_db"),
		Port:       getEnv("PORT", "8081"),
	}
}

// getEnv —  берет из env или возвращает дефолт
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func checkDBConnection(cfg Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
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

func init() {
	// Только загрузка файла .env. Никаких дефолтов!
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️ .env не найден или ошибка чтения: %v", err)
	}
}

func main() {
	//  Получаем готовую конфигурацию
	//  Выводим на экран
	cfg := loadConfig()
	log.Printf("📋 Конфигурация: host=%s, port=%s, db=%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	//  Проверяем подключение к БД
	db, err := checkDBConnection(cfg)
	if err != nil {
		log.Fatalf("❌ Критическая ошибка: не удалось подключиться к БД: %v", err)
	}
	defer db.Close()

	//  Настраиваем HTTP-сервер
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("EffectiveMobile API is running"))
	})

	addr := ":" + cfg.Port
	log.Printf("🚀 Сервер запущен на порту :%s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

}
