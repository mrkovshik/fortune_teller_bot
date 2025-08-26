package local

const (
	DorianGreyTitle     = "Оскар Уайлд - Портрет Дориана Грея"
	ThreeOnTheBoatTitle = "Дж.К.Джером - Трое в лодке, не считая собаки"
	GospodaGolovlevy    = "М.Е.Салтыков-Щедрин - Господа Головлёвы"
	DetiKapitanaGranta  = "Ж.Верн - Дети капитана Гранта"
	ZovKtulchu          = "Г.Лавкрафт - Зов Ктулху"
	ZoshenkoBest        = "М.Зощенко - Избранное"
)

var TitleToFileName = map[string]string{
	DorianGreyTitle:     "2.fb2",
	ThreeOnTheBoatTitle: "1.fb2",
	GospodaGolovlevy:    "3.fb2",
	DetiKapitanaGranta:  "4.epub",
	ZovKtulchu:          "5.epub",
	ZoshenkoBest:        "6.epub",
}

var FileNameToTitle = make(map[string]string)

func init() {
	FileNameToTitle = reverseMap(TitleToFileName)
}

func reverseMap(in map[string]string) (out map[string]string) {
	out = make(map[string]string)
	for k, v := range in {
		out[v] = k
	}
	return
}

func GetRandomBookTitle() string {
	for title := range TitleToFileName { // no need to use random - map iteration is random itself
		return title
	}
	return ""
}
