// main.go

package main

import (
	"./modules"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func main() {

	// Режим работы gin - на продакшене делать "ReleaseMode"
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

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
