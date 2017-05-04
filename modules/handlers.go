package modules

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/tebeka/selenium"
	"time"
)

func showIndexPage(context *gin.Context)  {

	// Обработка шаблона
	context.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title": "Главная страница",
		},
	)
}

func scrapVhfdx(context *gin.Context)  {



	var webDriver selenium.WebDriver
	var err error

	caps := selenium.Capabilities{
		"browserName":            "firefox",
		"webdriver.gecko.driver": "/usr/local/bin/geckodriver",
	}
	//
	//caps := selenium.Capabilities{
	//	"browserName":           "phantomjs",
	//	"phantomjs.binary.path": "/usr/local/bin/phantomjs",
	//}

	webDriver, err  = selenium.NewRemote(caps, "")

	if err != nil {
		panic(err)
	}
	defer webDriver.Quit()

	webDriver.MaximizeWindow("")

	// На сайт
	webDriver.Get("http://www.vhfdx.ru")
	time.Sleep(2 * time.Second)

	// Переходим на страницу Полевого дня.
	btn, _ := webDriver.FindElement(selenium.ByXPATH, "//a[@href='http://www.vhfdx.ru/field_day']")
	btn.Click()

	time.Sleep(2 * time.Second)

	// Переходим на страницу с таблицей участников.
	btn, _ = webDriver.FindElement(selenium.ByXPATH, "//a[contains(text(),'Таблица участников')]")
	btn.Click()
	time.Sleep(2 * time.Second)

	// Проверяем отображение страницы с таблицей участников.
	btn, _ = webDriver.FindElement(selenium.ByXPATH, "//title[contains(text(),'Российский УКВ портал - ПД  2016')]")
	btn.Click()
	time.Sleep(2 * time.Second)

	// Отсортировать по позывному
	btn, _ = webDriver.FindElement(selenium.ByXPATH, "//th/a[@id='farbikOrder_jos_fabrik_formdata_30.call']")
	btn.Click()
	time.Sleep(2 * time.Second)


	// Все строки с позывными со всех страниц в слайс
	overallResult := make([]string, 0, 200)
	for {
		// Считать данные с одной страницы
		var strings []string = readDateFromPage(webDriver)

		// Добавить к общему результату всех страниц
		for _, str := range strings {
			if str != "" {
				//fmt.Println(str)
				overallResult = append(overallResult, str)
			} else {
				break
			}
		}

		// Перейти на следующую страницу если есть ссылка
		var _, err = webDriver.FindElement(selenium.ByXPATH, "//a[@class='pagenav' and @title='Следующая']")

		if err == nil {
			// Переходим на следующую страницу.
			btn, _ = webDriver.FindElement(selenium.ByXPATH, "//a[@class='pagenav' and @title='Следующая']")
			btn.Click()
		} else {
			break
		}
	}



	var title, _ = webDriver.Title()
	//time.Sleep(10 * time.Second)

	// Закрыть браузер
	webDriver.Quit()


	context.HTML(
		http.StatusOK,
		"vhfdx.html",
		gin.H{
			"title": title,
			"callStrings": overallResult,
		},
	)
}

// Считывает данные участников с одной страницы сайта
func readDateFromPage(webDriver selenium.WebDriver) []string{

	// Позывные
	elements := make([]selenium.WebElement, 100)

	elements, _ = webDriver.FindElements(
		selenium.ByXPATH,
		"//tr[@class='oddrow0 fabrik_row ' or @class='oddrow1 fabrik_row ']/td[contains(@class,'call')]",
	)

	result := make([]string, 100)

	for index, item := range elements {
		result[index], _ = item.Text()
	}

	return result
}
