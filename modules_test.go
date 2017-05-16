package main

import (
	"./modules"
	"math"
	"testing"
)


// Тест для QRBFromDegrees()
func TestLatLongToQRA(t *testing.T) {

	type TestQrbValues struct {
		expectedQRB			int
		latitude1			float64
		longitude1			float64
		latitude2			float64
		longitude2			float64
	}

	var testQrbSlice = make([]TestQrbValues, 0 , 6)

	testQrbSlice = append(testQrbSlice,
		TestQrbValues{expectedQRB: 241, latitude1: 55.64, longitude1: 37.54, latitude2: 53.77, longitude2: 39.46},
		TestQrbValues{expectedQRB: 241, latitude1: 53.77, longitude1: 39.46, latitude2: 55.64, longitude2: 37.54},
		TestQrbValues{2518, 55.64, 37.54, 52.89, -1.46},
		TestQrbValues{1032, 53.77, 39.46, 53.40, 55.12},
		TestQrbValues{342, 53.77, 39.46, 56.06, 43.04},
		TestQrbValues{0, 53.77, 39.46, 53.77, 39.46},
	)

	// Первые две точки в тестовых данных одинаковые, просто поменяны местами - расстояния должны быть равны
	QRB1 := modules.QRBFromDegrees(testQrbSlice[0].latitude1, testQrbSlice[0].longitude1,
		testQrbSlice[0].latitude2, testQrbSlice[0].longitude2)
	QRB2 := modules.QRBFromDegrees(testQrbSlice[1].latitude1, testQrbSlice[1].longitude1,
		testQrbSlice[1].latitude2, testQrbSlice[1].longitude2)

	if QRB1 != QRB2 {
		t.Errorf("Расстояния не равны при перестановке местами точек: %v и %v", QRB1, QRB2)
	}


	// Пробежать по всем данным
	for _, testQrbItem := range testQrbSlice {

		actualQRB := modules.QRBFromDegrees(
			testQrbItem.latitude1,
			testQrbItem.longitude1,
			testQrbItem.latitude2,
			testQrbItem.longitude2)

		if actualQRB != testQrbItem.expectedQRB {
			t.Errorf("Неверное QRB: рассчитанное - %vкм, ожидаемое - %vкм", actualQRB, testQrbItem.expectedQRB)
		}
	}
}



// Тест для QraToLatLong()
func TestQraToLatLong(t *testing.T) {

	type TestQraValues struct {
		QRA             string
		expectedLat		float64
		expectedLong	float64
	}


	var testQraSlice = make([]TestQraValues, 0, 6)

	testQraSlice =  append(testQraSlice,
		TestQraValues{QRA: "KO85sp", expectedLat: 55.64, expectedLong: 37.54},
		TestQraValues{"IO92gv", 52.89, -1.46},
		TestQraValues{"AP98da", 68.02, -161.71},
		TestQraValues{"FJ00sa", 0.02, -78.46},
		TestQraValues{"FD55te", -54.81, -68.37},
		TestQraValues{"RB32ga", -77.98, 166.54},
		TestQraValues{"MR11op", 81.65, 63.21},
	)

	for _, testQraItem := range testQraSlice {

		precision := 1	// Число знаков после запятой для сравнения рассчитанных координат с ожидаемыми
		actualResultLatitude, actualResultLongitude := modules.QraToLatLong(testQraItem.QRA)


		if RoundTo(testQraItem.expectedLat, precision) != RoundTo(actualResultLatitude, precision) ||
				RoundTo(testQraItem.expectedLong, precision) != RoundTo(actualResultLongitude, precision) {
			t.Errorf("Неверные координаты квадрата %s: ожидалось %.2f и %.2f, получили %.2f и %.2f.",
				testQraItem.QRA,
				testQraItem.expectedLat, testQraItem.expectedLong,
				actualResultLatitude, actualResultLongitude,
			)
		}
	}
}

// Округляет к ближайшему с заданным количеством знаков после запятой
// value - округляемое значение
// digits - количество знаков после запятой
func RoundTo(value float64, digits int) float64 {
	var temp float64 = math.Pow10(digits)
	return (math.Floor(value * temp) + 0.5) / temp
}