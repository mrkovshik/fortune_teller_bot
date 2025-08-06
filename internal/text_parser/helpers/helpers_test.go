package helpers_test

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/text_parser/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RemoveTagsFromString", func() {
	DescribeTable("removes XML/HTML tags from string",
		func(input string, expected string) {
			result := helpers.RemoveTagsFromString(input)
			Expect(result).To(Equal(expected))
		},
		Entry("empty string", "", ""),
		Entry("string without tags", "plain text", "plain text"),
		Entry("simple tag", "<p>hello</p>", "hello"),
		Entry("nested tags", "<div><p>text</p></div>", "text"),
		Entry("self-closing tag", "text<br/>more", "textmore"),
		Entry("mixed content", "a <b>bold</b> move", "a bold move"),
		Entry("tags with attributes", `<a href="...">link</a>`, "link"),
		Entry("multiple tags", "<p>1</p><p>2</p>", "12"),
	)
})
