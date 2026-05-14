package main

import (
	"log"
	_ "subscriptions-service/docs"
	"subscriptions-service/internal/app"
)

// @title Subscription Service API
// @version 1.0
// @description REST API для агрегации данных об онлайн подписках пользователей
// @host localhost:8080
// @BasePath /
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
