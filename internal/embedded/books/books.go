package books

import (
	"embed"
)

//go:embed data/*
var booksFS embed.FS

func GetBooksFS() embed.FS {
	return booksFS
}
