package helpers_test

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/textparser/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CleanHTMLContent", func() {
	It("should remove script, style, and head tags with content", func() {
		html := `<html><head><title>Title</title></head><body><script>alert('x');</script><style>body{}</style><p>Hello</p></body></html>`
		result := helpers.CleanHTMLContent(html)
		Expect(result).To(Equal("Hello"))
	})

	It("should strip img tags", func() {
		html := `<p>Text before<img src="image.jpg">text after</p>`
		result := helpers.CleanHTMLContent(html)
		Expect(result).To(Equal("Text before text after"))
	})

	It("should preserve link text and remove the tag", func() {
		html := `<p>Read more <a href="https://example.com">here</a>.</p>`
		result := helpers.CleanHTMLContent(html)
		Expect(result).To(Equal("Read more here."))
	})

	It("should remove garbage like images/ and ch2.txt", func() {
		html := `<p>Chapter is in images/ch1.jpg and ch2.txt</p>`
		result := helpers.CleanHTMLContent(html)
		Expect(result).To(Equal("Chapter is in and"))
	})

	It("should normalize multiple spaces", func() {
		html := `<p>   Hello   world   </p>`
		result := helpers.CleanHTMLContent(html)
		Expect(result).To(Equal("Hello world"))
	})

	It("should clean complex HTML to plain text", func() {
		html := `<html><head><title>Ignored</title></head><body><h1>Main</h1><p>Some <a href="#">link</a> and <img src="bad.jpg"> image.</p><script>bad()</script></body></html>`
		result := helpers.CleanHTMLContent(html)
		Expect(result).To(Equal("Main Some link and image."))
	})
})
