package openai

import (
	"regexp"
	"strings"
)

var regexReplaceConsecutiveNewLine = regexp.MustCompile(`\n{2,}`)
var regexReplaceConsecutiveSpaces = regexp.MustCompile(`\s{2,}`)

func FormatPrompt(p string) string {
	return regexReplaceConsecutiveSpaces.ReplaceAllString(regexReplaceConsecutiveNewLine.ReplaceAllString(strings.ReplaceAll(strings.ReplaceAll(p, "\r", "\n"), "\t", "\n"), "\n"), " ")
}
