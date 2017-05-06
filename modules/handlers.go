package modules

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/tebeka/selenium"
	"time"
	"strings"
	"fmt"
	"strconv"
	"../models"
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

	var browser string
	browser = "phantom"		// Закомментировать для запуска chrome
	var caps selenium.Capabilities

	if browser == "phantom" {
		caps = selenium.Capabilities{
			"browserName":           "phantomjs",
			"phantomjs.binary.path": "/usr/local/bin/phantomjs",
		}
	} else {
		caps = selenium.Capabilities{
			"browserName":            "firefox",
			"webdriver.gecko.driver": "/usr/local/bin/geckodriver",
		}
	}


	webDriver, err  = selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}
	defer webDriver.Quit()

	webDriver.MaximizeWindow("")

	fmt.Println("Переход на страницу участников")
	err = webDriver.Get("http://www.vhfdx.ru/component/option,com_fabrik/Itemid,307/")
	time.Sleep(5 * time.Second)
	if err != nil {
		panic(err)
	}

	// Проверяем отображение страницы с таблицей участников.
	fmt.Println("Проверяем отображение страницы с таблицей участников")
	//btn, _ := webDriver.FindElement(selenium.ByXPATH, "//title[contains(text(),'Российский УКВ портал - ПД  2016')]")
	btn, err := webDriver.FindElement(selenium.ByXPATH, "//table[@class='adminlist']")
	if err != nil {
		panic(err)
	}
	time.Sleep(3 * time.Second)

	// Отсортировать по позывному
	fmt.Println("Отсортировать по позывному")
	btn, err = webDriver.FindElement(selenium.ByXPATH, "//th/a[@id='farbikOrder_jos_fabrik_formdata_30.call']")
	if err != nil {
		panic(err)
	}
	btn.Click()
	time.Sleep(3 * time.Second)


	// Все строки с позывными со всех страниц в слайс
	maxCountContestant := 200		// Максмальное количество участников соревнований, слайс растянется при необходимости
	overallResult := make([]string, 0, maxCountContestant)

	for {
		// Считать данные с одной страницы
		var contestantStrings []string = readDateFromPage(webDriver)

		// Добавить к общему результату всех страниц
		for _, str := range contestantStrings {
			if str != "" {
				//fmt.Println(str)
				overallResult = append(overallResult, str)
			} else {
				break
			}
		}

		// Перейти на следующую страницу если есть ссылка
		fmt.Println("Проверка наличия ссылки на Следующую страницу")
		var _, err = webDriver.FindElement(selenium.ByXPATH, "//a[@class='pagenav' and @title='Следующая']")
		time.Sleep(5 * time.Second)

		if err == nil {
			// Переходим на следующую страницу.
			fmt.Println("Переходим на Следующую страницу")
			btn, _ = webDriver.FindElement(selenium.ByXPATH, "//a[@class='pagenav' and @title='Следующая']")
			err = btn.Click()
			if err != nil {
				panic(err)
			}
			time.Sleep(7 * time.Second)

			// Проверяем отображение страницы с таблицей участников.
			count := 10
			for i:=0; i<count; i++ {
				fmt.Printf("Проверяем отображение страницы с таблицей участников. Попытка N%d\n", i+1)
				btn, err = webDriver.FindElement(selenium.ByXPATH, "//table[@class='adminlist']")
				if err == nil {
					break
				}
			}
			if err != nil {
				fmt.Println("Не отобразилась следующая страница.")
				panic(err)
			}

		} else {
			fmt.Println("Нет ссылки на Следующую страницу")
			break
		}
	}

	var title, _ = webDriver.Title()

	// Закрыть браузер
	webDriver.Quit()


	/// Обработка результатов
	// Генерация матрицы диапазонов
	overallResult = makeBandMatrix(overallResult)

	// Вычисление расстояния

	// Вычисление азимута

	// Вычисление обратного азимута

	// Получить всех участников
	contestantCount := len(overallResult)	// Количесво участников
	contestantList := make([]models.Contestant, contestantCount)
	toHTMLTable(overallResult, contestantList)

	context.HTML(
		http.StatusOK,
		"vhfdx.html",
		gin.H{
			"title": title,
			"contestantCount": contestantCount,
			"contestant": contestantList,
		},
	)
}


// Внести данные об участнике в структуру для контекста темплейта
func toHTMLTable(allContestant []string, contestantList []models.Contestant) {
	for idx, str := range allContestant {
		contestantList[idx].ID, _ = strconv.Atoi(strings.Split(str, "||")[0])
		contestantList[idx].Call = strings.Split(str, "||")[1]
		contestantList[idx].QRA = strings.Split(str, "||")[2]
		contestantList[idx].Band_2m = strings.Split(str, "||")[3]
		contestantList[idx].Band_70cm = strings.Split(str, "||")[4]
		contestantList[idx].Band_23cm = strings.Split(str, "||")[5]
		contestantList[idx].Band_5cm = strings.Split(str, "||")[6]
		contestantList[idx].Band_3cm = strings.Split(str, "||")[7]
		contestantList[idx].Band_1cm = strings.Split(str, "||")[8]
		contestantList[idx].Info = strings.Split(str, "||")[9]
	}
}


// Создаёт матрицу диапазонов
func makeBandMatrix(allContestant []string) []string {

	maxCountContestant := 200	// Максмальное количество участников соревнований, слайс растянется при необходимости
	outputStringWithBandMatrix := make([]string, 0, maxCountContestant)

	for idx, str := range allContestant {
		bandString := ""

		// Парсим диапазоны
		bands := strings.Split(str, "||")[2]

		info := strings.Split(str, "||")[3]
		if info == "" {							// Если ячейка с 'Дополнительной информацией' пуста
			info = " "
		}

		if strings.Contains(bands, "144МГц") {
			bandString += "2м"
		} else {
			bandString += " "
		}

		if strings.Contains(bands, "432МГц") {
			bandString += "||70см"
		} else {
			bandString += "|| "
		}

		if strings.Contains(bands, "1296МГц") {
			bandString += "||23см"
		} else {
			bandString += "|| "
		}

		if strings.Contains(bands, "5.7ГГц") {
			bandString += "||5см"
		} else {
			bandString += "|| "
		}

		if strings.Contains(bands, "10ГГц") {
			bandString += "||3м"
		} else {
			bandString += "|| "
		}

		if strings.Contains(bands, "24ГГц") {
			bandString += "||1.2см"
		} else {
			bandString += "|| "
		}

		outputStringWithBandMatrix = append(outputStringWithBandMatrix,
			strconv.Itoa(idx + 1) + "||" +
			strings.Split(str, "||")[0] + "||" +
			strings.Split(str, "||")[1] + "||" +
			bandString + "||" +
			info)
	}
	return outputStringWithBandMatrix
}


// Считывает данные участников с одной страницы сайта
func readDateFromPage(webDriver selenium.WebDriver) []string{

	callsOnPage := 100		// Максимальное лоличество позывных на одной страниц
	var err error

	// Позывные
	callElements := make([]selenium.WebElement, callsOnPage)
	fmt.Println("Считываем позывные")
	callElements, err = webDriver.FindElements(selenium.ByXPATH,
		"//tr[@class='oddrow0 fabrik_row ' or @class='oddrow1 fabrik_row ']/td[contains(@class,'call')]",)
	if err != nil {
		panic(err)
	}

	// Квадраты
	qraElements := make([]selenium.WebElement, callsOnPage)
	fmt.Println("Считываем квадраты")
	qraElements, err = webDriver.FindElements(selenium.ByXPATH,
		"//tr[@class='oddrow0 fabrik_row ' or @class='oddrow1 fabrik_row ']/td[contains(@class,'qra')]")
	if err != nil {
		panic(err)
	}

	// Диапазоны
	bandsElements := make([]selenium.WebElement, callsOnPage)
	fmt.Println("Считываем диапазоны")
	bandsElements, err = webDriver.FindElements(selenium.ByXPATH,
		"//tr[@class='oddrow0 fabrik_row ' or @class='oddrow1 fabrik_row ']/td[contains(@class,'band')]")
	if err != nil {
		panic(err)
	}

	// Дополнительная информация
	infoElements := make([]selenium.WebElement, callsOnPage)
	fmt.Println("Считываем дополнительную информацию")
	infoElements, err = webDriver.FindElements(selenium.ByXPATH,
		"//tr[@class='oddrow0 fabrik_row ' or @class='oddrow1 fabrik_row ']/td[contains(@class,'info')]")
	if err != nil {
		panic(err)
	}


	result := make([]string, callsOnPage)

	for i := 0; i < len(callElements); i++ {

		call, _ := callElements[i].Text()
		qra, _ := qraElements[i].Text()
		bands, _ := bandsElements[i].Text()
		info, _ := infoElements[i].Text()

		result[i] = fmt.Sprintf("%s||%s||%s||%s", 
			strings.ToTitle(strings.Replace(call, " ", "", -1)),		// Всё в верхний регистр и убрать пробелы
			strings.ToTitle(strings.Replace(qra, " ", "", -1)), 
			bands, 
			info)
	}
	return result
}
