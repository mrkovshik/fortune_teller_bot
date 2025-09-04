package epub_test

import (
	"time"

	"github.com/mrkovshik/fortune_teller_bot/internal/textparser/epub"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var parser = epub.NewTextParser()

var _ = Describe("", func() {

	It("", func() {
		for _, book := range TestBooks {
			sent, err := parser.ParseRandomSentence(book, time.Now().UnixNano())
			Expect(err).To(Succeed())
			Expect(sent).NotTo(BeNil())
			break
		}
	})
})
