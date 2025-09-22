package main

import (
	"os"

	"github.com/4udiwe/subscription-service/docs"
	"github.com/4udiwe/subscription-service/internal/app"
)

// @title Subscriptions Service
// @version 1.0
// @description Сервис для агрегации онлайн-подписок (тестовое задание).

// @contact.email sharifkulov.work@gmail.com

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {

	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	app := app.New(os.Getenv("CONFIG_PATH"))
	app.Start()
}
