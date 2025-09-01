package postgres

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	booksmeta "github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testID = 1

var _ = Describe("Posgres DB storage", func() {
	It("Prepares list query", func() {
		query, args := prepareListQuery(booksmeta.WithLanguage(config.Russian), booksmeta.WithAuthor("test author"))
		Expect(query).To(Equal("SELECT id, title, author, language, format FROM books WHERE author ILIKE $1 AND language = $2"))
		Expect(args).To(ContainElements("%test author%", config.Russian))
	})
})
