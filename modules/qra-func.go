package modules

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

// Проверяет валидность введённого QRA
func checkQRA(qra string) error {

	qra = strings.ToUpper(qra)
	var result error = nil

	// Первая буква от A до R?
	match, _ := regexp.MatchString("[A-R]{2}\\d{2}[A-X]{2}", qra)
	if !match {
		result = fmt.Errorf("Введён неверный формат квадрата QRA")
	}
	return result
}


// Вычисляет прямой и обратный азимуты для двух квадратов
func AzimuthsFromQRA(qra1 string, qra2 string) (int, int) {
	lat1, long1 := QraToLatLong(qra1)
	lat2, long2 := QraToLatLong(qra2)
	azi, revAzi := Azimuths(lat1, long1, lat2, long2)
	return azi, revAzi
}

// Вычисляет прямой и обратный азимуты в градусах двух точек с координатами в градусах.
func Azimuths(lat1, long1, lat2, long2 float64) (int, int) {

	//var EarthRadius float64 = 6372795.0			// Радиус Земли в метрах

	// В радианы
	lat1 = DegreesToRadians(lat1)
	lat2 = DegreesToRadians(lat2)
	long1 = DegreesToRadians(long1)
	long2 = DegreesToRadians(long2)

	// Косинусы и синусы широт и разницы долгот
	cosLat1 := math.Cos(lat1)
	cosLat2 := math.Cos(lat2)
	sinLat1 := math.Sin(lat1)
	sinLat2 := math.Sin(lat2)
	longDelta := long2 - long1
	cosLongDelta := math.Cos(longDelta)
	sinLongDelta := math.Sin(longDelta)

	// Начальный азимут
	x := cosLat1 * sinLat2 - sinLat1 * cosLat2 * cosLongDelta
	y := sinLongDelta * cosLat2
	z := RadiansToDegrees(math.Atan(-y/x))

	if x < 0 {
		z = z + 180.0
	}
	var z2 float64 = float64(int(z + 180.0) % 360 - 180)
	z2 = -1 * DegreesToRadians(z2)

	var azimuth int = int(RadiansToDegrees(z2 - (2*math.Pi*math.Floor(z2/(2*math.Pi)))))
	var reverseBearing int
	if azimuth < 180 {
		reverseBearing = azimuth + 180
	} else {
		reverseBearing = azimuth - 180
	}

	return azimuth, reverseBearing
}


// Вычисляет QRB в километрах между двумя точками с координатами в градусах.
func QRBFromDegrees(lat1, long1, lat2, long2 float64) int {

	k1 := math.Cos(DegreesToRadians(lat1)) *
			math.Cos(DegreesToRadians(lat2)) *
			math.Cos(DegreesToRadians(math.Abs(long1-long2)))

	k2 := math.Sin(DegreesToRadians(lat1)) * math.Sin(DegreesToRadians(lat2))

	qrb := 111.25 * RadiansToDegrees(math.Acos(k1 + k2))

	return int(qrb)
}

// Вычисляет QRB в километрах между двумя QRA квадратами
func QRBFromQRA(qra1 string, qra2 string) int {

	lat1, long1 := QraToLatLong(qra1)
	lat2, long2 := QraToLatLong(qra2)
	qrb := QRBFromDegrees(lat1, long1, lat2, long2)
	return qrb
}

// Получает градусы, возвращает радианы
func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Получает радианы, возвращае градусы
func RadiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}



// Получает квадрат QRA, возвращает широту и долготу в градусах
func QraToLatLong(qra string) (float64, float64) {

	var lat float64  // Широта в градусах
	var long float64 // Долгота в градусах

	// Мапы для широты
	latMap1 := map[string]float64 {
		"A": -90, "B": -80, "C": -70, "D": -60, "E": -50, "F": -40, "G": -30, "H": -20, "I": -10,
		"J": 0, "K": 10, "L": 20, "M": 30, "N": 40, "O": 50, "P": 60, "Q": 70, "R": 80,
	}
	latMap2 := map[string]float64 {
		"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
	}
	latMap3 := map[string]float64 {
		"A": 0.0208, "B": 0.0625, "C": 0.1042, "D": 0.1458, "E": 0.1875, "F": 0.2292, "G": 0.2708, "H": 0.3125,
		"I": 0.3542, "J": 0.3958, "K": 0.4375, "L": 0.4792,	"M": 0.5208, "N": 0.5625, "O": 0.6042, "P": 0.6458,
		"Q": 0.6875, "R": 0.7292, "S": 0.7708, "T": 0.8125, "U": 0.8542, "V": 0.8958, "W": 0.9375, "X": 0.9792,
	}

	// Мапы для долготы
	longMap1 := map[string]float64 {
		"A": -180, "B": -160, "C": -140, "D": -120, "E": -100, "F": -80, "G": -60, "H": -40, "I": -20,
		"J": 0, "K": 20, "L": 40, "M": 60, "N": 80, "O": 100, "P": 120, "Q": 140, "R": 160,
	}
	longMap2 := map[string]float64 {
		"0": 0, "1": 2, "2": 4, "3": 6, "4": 8, "5": 10, "6": 12, "7": 14, "8": 16, "9": 18,
	}
	longMap3 := map[string]float64 {
		"A": 0.0417, "B": 0.1250, "C": 0.2083, "D": 0.2917, "E": 0.3750, "F": 0.4583, "G": 0.5417, "H": 0.6250,
		"I": 0.7083, "J": 0.7917, "K": 0.8750, "L": 0.9583,	"M": 1.0417, "N": 1.1250, "O": 1.2083, "P": 1.2917,
		"Q": 1.3750, "R": 1.4583, "S": 1.5417, "T": 1.6250, "U": 1.7083, "V": 1.7917, "W": 1.8750, "X": 1.9583,
	}

	lat = latMap1[strings.ToUpper(qra[1:2])] + latMap2[qra[3:4]] + latMap3[strings.ToUpper(qra[5:6])]
	long = longMap1[strings.ToUpper(qra[0:1])] + longMap2[qra[2:3]] + longMap3[strings.ToUpper(qra[4:5])]

	return lat, long
}
