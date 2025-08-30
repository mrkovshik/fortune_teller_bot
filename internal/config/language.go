package config

type Language string

const (
	English Language = "eng"
	Russian Language = "rus"
)

var SupportedLanguages = map[Language]string{
	English: "English",
	Russian: "Русский",
}

func (l Language) String() string {
	return string(l)
}
