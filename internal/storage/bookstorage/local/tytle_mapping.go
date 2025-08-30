package local

import "github.com/mrkovshik/fortune_teller_bot/internal/config"

const (
	DorianGreyTitle     = "Оскар Уайлд - Портрет Дориана Грея"
	ThreeOnTheBoatTitle = "Дж.К.Джером - Трое в лодке, не считая собаки"
	GospodaGolovlevy    = "М.Е.Салтыков-Щедрин - Господа Головлёвы"
	DetiKapitanaGranta  = "Ж.Верн - Дети капитана Гранта"
	ZovKtulchu          = "Г.Лавкрафт - Зов Ктулху"
	ZoshenkoBest        = "М.Зощенко - Избранное"
	MelvilleMobyDick    = "H.Melville - Moby Dick"
	Frankenstein        = "M.Shelley - Frankenstein"
)

var TitleToFileName = map[config.Language]map[string]string{
	config.English: {
		MelvilleMobyDick: "2.epub",
		Frankenstein:     "1.epub",
	},
	config.Russian: {
		DorianGreyTitle:     "2.fb2",
		ThreeOnTheBoatTitle: "1.fb2",
		GospodaGolovlevy:    "3.fb2",
		DetiKapitanaGranta:  "4.epub",
		ZovKtulchu:          "5.epub",
		ZoshenkoBest:        "6.epub",
	},
}

var FileNameToTitle = make(map[config.Language]map[string]string)

func init() {
	FileNameToTitle = reverseMap(TitleToFileName)
}

func reverseMap(in map[config.Language]map[string]string) (out map[config.Language]map[string]string) {
	out = make(map[config.Language]map[string]string)
	for lang, books := range in {
		out[lang] = make(map[string]string)
		for k, v := range books {
			out[lang][v] = k
		}
	}
	return
}

func GetRandomBookTitle(language config.Language) string {
	for title := range TitleToFileName[language] { // no need to use random - map iteration is random itself
		return title
	}
	return ""
}
