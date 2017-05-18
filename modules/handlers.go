package modules

import (
	"../models"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/tebeka/selenium"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	log.SetOutput(os.Stdout)				// В консоль
	log.SetLevel(log.InfoLevel)				// Уровень логирования
	log.SetFormatter(&log.TextFormatter{})	// Текстовые логи

	// Квадрат из формы
	myQRA := context.PostForm("my_qra")
	log.Infof("Получен myQRA: %s", myQRA)

	var err error

	// Проверка валидности введённого квадрата
	log.Info("Проверяем валидность введённого квадрата")
	err = CheckQRA(myQRA)
	if err != nil {
		context.HTML(
			http.StatusOK,
			"message.html",
			gin.H{
				"title": "Ошибка",
				"message1": fmt.Sprintf("%s: ", err),
				"message3": myQRA,
			},
		)
	} else {

		var webDriver selenium.WebDriver

		var browser string
		browser = "phantom" // Закомментировать для запуска chrome
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

		webDriver, err = selenium.NewRemote(caps, "")
		if err != nil {
			panic(err)
		}
		defer webDriver.Quit()

		webDriver.MaximizeWindow("")

		log.Info("Переход на страницу участников")
		err = webDriver.Get("http://www.vhfdx.ru/component/option,com_fabrik/Itemid,307/")
		time.Sleep(5 * time.Second)
		if err != nil {
			panic(err)
		}

		// Проверяем отображение страницы с таблицей участников.
		log.Info("Проверяем отображение страницы с таблицей участников")
		//btn, _ := webDriver.FindElement(selenium.ByXPATH, "//title[contains(text(),'Российский УКВ портал - ПД  2016')]")
		btn, err := webDriver.FindElement(selenium.ByXPATH, "//table[@class='adminlist']")
		if err != nil {
			panic(err)
		}
		time.Sleep(3 * time.Second)

		// Отсортировать по позывному
		log.Info("Отсортировать по позывному")
		btn, err = webDriver.FindElement(selenium.ByXPATH, "//th/a[@id='farbikOrder_jos_fabrik_formdata_30.call']")
		if err != nil {
			panic(err)
		}
		btn.Click()
		time.Sleep(3 * time.Second)

		// Все строки с позывными со всех страниц в слайс
		maxCountContestant := 200 // Максмальное количество участников соревнований, слайс растянется при необходимости
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
			log.Info("Проверка наличия ссылки на Следующую страницу")
			var _, err= webDriver.FindElement(selenium.ByXPATH, "//a[@class='pagenav' and @title='Следующая']")
			time.Sleep(5 * time.Second)

			if err == nil {
				// Переходим на следующую страницу.
				log.Info("Переходим на Следующую страницу")
				btn, _ = webDriver.FindElement(selenium.ByXPATH, "//a[@class='pagenav' and @title='Следующая']")
				err = btn.Click()
				if err != nil {
					panic(err)
				}
				time.Sleep(7 * time.Second)

				// Проверяем отображение страницы с таблицей участников.
				count := 5
				for i := 0; i < count; i++ {
					log.Infof("Проверяем отображение страницы с таблицей участников. Попытка N%d", i+1)
					btn, err = webDriver.FindElement(selenium.ByXPATH, "//table[@class='adminlist']")
					time.Sleep(5 * time.Second)
					if err == nil {
						break
					}
				}
				if err != nil {
					log.Info("Не отобразилась следующая страница.")
					panic(err)
				}

			} else {
				log.Info("Нет ссылки на Следующую страницу")
				break
			}
		}

		var title, _= webDriver.Title()

		// Закрыть браузер
		webDriver.Quit()

		/// Обработка результатов
		// Генерация матрицы диапазонов
		overallResult = makeBandMatrix(overallResult)

		// Получить всех участников
		contestantCount := len(overallResult) // Количесво участников
		contestantList := make([]models.Contestant, contestantCount)

		// Заполнить HTML таблицу информацией об участниках, QRB и азимутами
		toHTMLTable(overallResult, contestantList, myQRA)

		context.HTML(
			http.StatusOK,
			"vhfdx.html",
			gin.H{
				"title":           title,
				"myQRA":           myQRA,
				"contestantCount": contestantCount,
				"contestant":      contestantList,
			},
		)
	}
}



// Внести данные об участнике в структуру для контекста темплейта
func toHTMLTable(allContestant []string, contestantList []models.Contestant, myQRA string) {
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
		contestantList[idx].QRB = QRBFromQRA(myQRA, strings.Split(str, "||")[2])
		contestantList[idx].Azi, contestantList[idx].ReversAzi = AzimuthsFromQRA(myQRA, strings.Split(str, "||")[2])
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
			bandString += "||3см"
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
	log.Info("Считываем позывные")
	callElements, err = webDriver.FindElements(selenium.ByXPATH,
		"//tr[@class='oddrow0 fabrik_row ' or @class='oddrow1 fabrik_row ']/td[contains(@class,'call')]",)
	if err != nil {
		panic(err)
	}

	// Квадраты
	qraElements := make([]selenium.WebElement, callsOnPage)
	log.Info("Считываем квадраты")
	qraElements, err = webDriver.FindElements(selenium.ByXPATH,
		"//tr[@class='oddrow0 fabrik_row ' or @class='oddrow1 fabrik_row ']/td[contains(@class,'qra')]")
	if err != nil {
		panic(err)
	}

	// Диапазоны
	bandsElements := make([]selenium.WebElement, callsOnPage)
	log.Info("Считываем диапазоны")
	bandsElements, err = webDriver.FindElements(selenium.ByXPATH,
		"//tr[@class='oddrow0 fabrik_row ' or @class='oddrow1 fabrik_row ']/td[contains(@class,'band')]")
	if err != nil {
		panic(err)
	}

	// Дополнительная информация
	infoElements := make([]selenium.WebElement, callsOnPage)
	log.Info("Считываем дополнительную информацию")
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
