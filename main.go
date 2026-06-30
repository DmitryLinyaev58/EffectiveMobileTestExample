package main

import (
	"log"

	"github.com/DmitryLinyaev58/EffectiveMobileTestExample/bootstrap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	log.Println("...пуск приложения...")
	if err := bootstrap.Run(); err != nil {
		log.Fatalf("❌ Критическая ошибка при запуске приложения: %v", err)
	}
}
