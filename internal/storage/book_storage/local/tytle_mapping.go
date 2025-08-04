package local

const (
	DorianGreyTitle     = "Оскар Уайлд - Портрет Дориана Грея"
	ThreeOnTheBoatTitle = "Дж.К.Джером - Трое в лодке, не считая собаки"
	GospodaGolovlevy    = "М.Е.Салтыков-Щедрин - Господа Головлёвы"
	ZolotoyTelenok      = "И.Ильф и Е.Петров - Золотой Телёнок"
)

var TitleToFileName = map[string]string{
	DorianGreyTitle:     "1",
	ThreeOnTheBoatTitle: "2",
	GospodaGolovlevy:    "3",
	//ZolotoyTelenok:      "4",
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
