package embedded

import (
	"embed"
)

//go:embed books/*.fb2
var booksFS embed.FS

func GetBooksFS() embed.FS {
	return booksFS
}
