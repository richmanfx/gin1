package main

import (
	"testing"
	"./modules"
)

// Тест для QraToLatLong()
func TestQraToLatLong(t *testing.T)  {

	type TestQraValues struct {
		QRA             string
		expectedLat		int
		expectedLong	int
	}

	var testQraSlice = make([]TestQraValues, 0, 6)

	testQraSlice =  append(testQraSlice,
		TestQraValues{QRA: "KO85sp", expectedLat: 56, expectedLong: 38},
		TestQraValues{"IO92gv", 53, -1},
		TestQraValues{"AP98da", 68, -162},
		TestQraValues{"FJ00SA", 0, -78},
		TestQraValues{"FD55te", -55, -68},
		TestQraValues{"RB32ga", -78, 167},
		TestQraValues{"MR11op", 82, 63},
	)

	for _, testQraItem := range testQraSlice {

		actualResultLatitude, actualResultLongitude := modules.QraToLatLong(testQraItem.QRA)

		if (testQraItem.expectedLat != actualResultLatitude) ||
		   (testQraItem.expectedLong != actualResultLongitude) {
			t.Fatalf("Неверные координаты квадрата %s: ожидалось %d и %d, получили %d и %d.",
				testQraItem.QRA,
				testQraItem.expectedLat, testQraItem.expectedLong,
				actualResultLatitude, actualResultLongitude,
			)
		}
	}
}

