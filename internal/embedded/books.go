package embedded

import (
	"embed"
)

//go:embed books/*
var booksFS embed.FS

func GetBooksFS() embed.FS {
	return booksFS
}
