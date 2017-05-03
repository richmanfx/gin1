package modules

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/tebeka/selenium"
	"time"
)

func showIndexPage(context *gin.Context)  {

	// Обработка шаблона
	context.HTML(http.StatusOK, "index.html", gin.H{
													"title": "Главная страница",
											  },
	)
}

func scrapVhfdx(context *gin.Context)  {



	var webDriver selenium.WebDriver
	var err error

	//caps := selenium.Capabilities{
	//	"browserName":            "firefox",
	//	"webdriver.gecko.driver": "/usr/local/bin/geckodriver",
	//}

	caps := selenium.Capabilities{
		"browserName":           "phantomjs",
		"phantomjs.binary.path": "/usr/local/bin/phantomjs",
	}

	webDriver, err  = selenium.NewRemote(caps, "")

	if err != nil {
		panic(err)
	}
	defer webDriver.Quit()

	webDriver.MaximizeWindow("")

	// Get simple playground interface
	webDriver.Get("http://www.vhfdx.ru")

	// Переходим на страницу Полевого дня.
	btn, _ := webDriver.FindElement(selenium.ByXPATH, "//a[@href='http://www.vhfdx.ru/field_day']")
	btn.Click()

	time.Sleep(3 * time.Second)

	// Переходим на страницу с таблицей участников.
	btn, _ = webDriver.FindElement(selenium.ByXPATH, "//a[contains(text(),'Таблица участников')]")
	btn.Click()
	time.Sleep(3 * time.Second)

	// Проверяем отображение страницы с таблицей участников.
	btn, _ = webDriver.FindElement(selenium.ByXPATH, "//title[contains(text(),'Российский УКВ портал - ПД  2016')]")
	btn.Click()
	time.Sleep(3 * time.Second)

	// Отсортировать по позывному
	btn, _ = webDriver.FindElement(selenium.ByXPATH, "//th/a[@id='farbikOrder_jos_fabrik_formdata_30.call']")
	btn.Click()
	time.Sleep(3 * time.Second)

	// Переходим на следующую страницу.
	btn, _ = webDriver.FindElement(selenium.ByXPATH, "//a[@class='pagenav' and @title='Следующая']")
	btn.Click()


	var title, _ = webDriver.Title()
	time.Sleep(10 * time.Second)


	context.HTML(http.StatusOK, "index.html", gin.H{
		"title": title,
		"fromVhfdx": "Большой привет. Смотри на title: " + title,
	},
	)
}
