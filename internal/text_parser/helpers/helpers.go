package helpers

import (
	"regexp"
)

// RemoveTagsFromString removes xml tags from string
func RemoveTagsFromString(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(s, "")
}
