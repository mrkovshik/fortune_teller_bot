package local

const (
	DorianGreyTitle     = "Оскар Уайлд - Портрет Дориана Грея"
	ThreeOnTheBoatTitle = "Дж. К. Джером - Трое в лодке, не считая собаки"
)

var titleToFileName map[string]string = map[string]string{
	DorianGreyTitle:     "oskar_wild_dorian_grey.fb2",
	ThreeOnTheBoatTitle: "jerom_troe_v_lodke.fb2",
}

var fileNameToTitle = make(map[string]string)

func init() {
	fileNameToTitle = reverseMap(titleToFileName)
}

func reverseMap(in map[string]string) (out map[string]string) {
	out = make(map[string]string)
	for k, v := range in {
		out[v] = k
	}
	return
}
