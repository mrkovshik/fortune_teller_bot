package local

var titleToFileName map[string]string = map[string]string{
	"Оскар Уайлд - Портрет Дориана Грея": "oskar_wild_dorian_grey.fb2",
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
