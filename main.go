// main.go

package main

import (
	"github.com/gin-gonic/gin"
	"./modules"
)

var Router *gin.Engine

func main() {
	// Роутер по-умолчанию в Gin
	Router = gin.Default()

	// Загрузить шаблоны
	Router.LoadHTMLGlob("templates/*")

	// Загрузить статику
	Router.Static("/css", "css")
	Router.Static("/js", "js")

	// Проинитить роуты
	modules.InitRoutes(Router)

	// Запустить приложение
	Router.Run(":8282")
}
