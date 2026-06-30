package main

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/DmitryLinyaev58/EffectiveMobileTestExample/bootstrap"
)

func main() {

	log.Println("...пуск приложения...")
	if err := bootstrap.Run(); err != nil {
		log.Fatalf("❌ Критическая ошибка при запуске приложения: %v", err)
	}

}
