package modules

import "strings"


// Получает квадрат QRA, возвращает широту и долготу в градусах
func QraToLatLong(qra string) (int, int) {

	var lat int  // Широта в градусах
	var long int // Долгота в градусах

	// Мапы для широты
	latMap1 := map[string]int {
		"A": -90, "B": -80, "C": -70, "D": -60, "E": -50, "F": -40, "G": -30, "H": -20, "I": -10,
		"J": 0, "K": 10, "L": 20, "M": 30, "N": 40, "O": 50, "P": 60, "Q": 70, "R": 80,
	}
	latMap2 := map[string]int {
		"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
	}
	latMap3 := map[string]int {
		"A": 0, "B": 0, "C": 0, "D": 0, "E": 0, "F": 0, "G": 0, "H": 0, "I": 0,	"J": 0, "K": 0, "L": 0,
		"M": 1, "N": 1, "O": 1, "P": 1, "Q": 1, "R": 1, "S": 1, "T": 1, "U": 1,	"V": 1, "W": 1, "X": 1,
	}

	// Мапы для долготы
	longMap1 := map[string]int {
		"A": -180, "B": -160, "C": -140, "D": -120, "E": -100, "F": -80, "G": -60, "H": -40, "I": -20,
		"J": 0, "K": 20, "L": 40, "M": 60, "N": 80, "O": 100, "P": 120, "Q": 140, "R": 160,
	}
	longMap2 := map[string]int {
		"0": 0, "1": 2, "2": 4, "3": 6, "4": 8, "5": 10, "6": 12, "7": 14, "8": 16, "9": 18,
	}
	longMap3 := map[string]int {
		"A": 0, "B": 0, "C": 0, "D": 0, "E": 0, "F": 0, "G": 1, "H": 1, "I": 1,	"J": 1, "K": 1, "L": 1,
		"M": 1, "N": 1, "O": 1, "P": 1, "Q": 1, "R": 1, "S": 2, "T": 2, "U": 2,	"V": 2, "W": 2, "X": 2,
	}

	lat = latMap1[strings.ToUpper(qra[1:2])] + latMap2[qra[3:4]] + latMap3[strings.ToUpper(qra[5:6])]
	long = longMap1[strings.ToUpper(qra[0:1])] + longMap2[qra[2:3]] + longMap3[strings.ToUpper(qra[4:5])]

	return lat, long
}
