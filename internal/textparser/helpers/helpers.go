package helpers

import (
	"regexp"
	"strings"
)

func CleanHTMLContent(input string) string {
	// Remove <script>...</script>, <style>...</style>, <head>...</head> blocks
	reScript := regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	input = reScript.ReplaceAllString(input, " ")

	reStyle := regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	input = reStyle.ReplaceAllString(input, " ")

	reHead := regexp.MustCompile(`(?is)<head[^>]*>.*?</head>`)
	input = reHead.ReplaceAllString(input, "")

	// Remove all <img ...> tags
	reImg := regexp.MustCompile(`(?is)<img[^>]*>`)
	input = reImg.ReplaceAllString(input, " ")

	// Remove <a ...> tags but keep the inner text
	reA := regexp.MustCompile(`(?is)<a[^>]*>(.*?)</a>`)
	input = reA.ReplaceAllString(input, "$1")

	// Remove all remaining HTML tags
	reTags := regexp.MustCompile(`(?is)<[^>]+>`)
	input = reTags.ReplaceAllString(input, " ")

	// Remove leftover HTML attributes like href and src
	reAttrs := regexp.MustCompile(`(?i)\s*(href|src)="[^"]*"`)
	input = reAttrs.ReplaceAllString(input, " ")

	// Remove garbage like "images/..." or "ch2."
	reGarbage := regexp.MustCompile(`(?i)(images/[^ ]+|ch\d+\.\w*)`)
	input = reGarbage.ReplaceAllString(input, " ")

	// Replace multiple spaces with a single space
	reSpaces := regexp.MustCompile(`\s+`)
	input = reSpaces.ReplaceAllString(input, " ")

	input = strings.ReplaceAll(input, "  ", " ")

	return strings.TrimSpace(input)
}
