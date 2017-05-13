package main

import (
	"./modules"
	"math"
	"testing"
)

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