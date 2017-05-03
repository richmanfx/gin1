// routes.go

package modules

import "github.com/gin-gonic/gin"

func InitRoutes(router *gin.Engine) {

	// Роутинг главной страницы: метод, путь -> обработчик

	router.Handle("GET", "/", showIndexPage)
	router.Handle("GET","vhfdx", scrapVhfdx)

}
