package messages

import (
	"regexp"
	"strings"
)

var codeFenceOpen = regexp.MustCompile("(?m)^```[[:alnum:]_+-]*\\s*")
var markdownMarks = regexp.MustCompile("[*_`#>\\[\\]]+")
var whitespace = regexp.MustCompile("\\s+")

func PlainTextFromMarkdown(input string) string {
	stripped := codeFenceOpen.ReplaceAllString(input, "")
	stripped = markdownMarks.ReplaceAllString(stripped, "")
	return strings.TrimSpace(whitespace.ReplaceAllString(stripped, " "))
}
