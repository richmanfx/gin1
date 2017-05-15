package models

type Contestant struct {
	// Названия с заглавных букв - доступны во всём пакете
	ID			int
	Call		string
	QRA			string
	Band_2m		string
	Band_70cm	string
	Band_23cm	string
	Band_5cm	string
	Band_3cm	string
	Band_1cm	string
	QRB			int
	Azi			int
	ReversAzi	int
	Info		string
}
