// routes.go

package modules

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	// Роутинг страницы: метод, путь -> обработчик
	router.Handle("GET", "/go", showIndexPage)
	router.Handle("POST", "/go/fd", fieldDayContestant)
	router.Handle("POST", "/go/russia-cup", RussiaСhampionshipContestant)
}
