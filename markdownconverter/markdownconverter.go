package markdownconverter

import (
	"errors"
	"regexp"
	"strings"
)

var (
	TagMap = map[*regexp.Regexp]string{
		regexp.MustCompile(`\*\*(.*?)\*\*`): "<b>$1</b>",
		regexp.MustCompile(`\b_(.*?)_\b`):   "<i>$1</i>",
		regexp.MustCompile("`([^`]+)`"):     "<tt>$1</tt>",
	}

	preformattedBlockOpeningTag = "<pre>\n"
	preformattedBlockClosingTag = "</pre>\n"
	backtick                    = "```"
	paragraphOpeningTag         = "<p>"
	paragraphClosingTag         = "</p>"
)

func processPreformattedBlock(result *strings.Builder, isPreformattedBlock *bool) {
	*isPreformattedBlock = !*isPreformattedBlock
	if *isPreformattedBlock {
		result.WriteString(preformattedBlockOpeningTag)
	} else {
		result.WriteString(preformattedBlockClosingTag)
	}
}

func processParagraph(result *strings.Builder, isParagraphOpen *bool, line string) {
	trimmedLine := strings.TrimSpace(line)
	if trimmedLine == "" {
		if *isParagraphOpen {
			result.WriteString(paragraphClosingTag + "\n")
			*isParagraphOpen = false
		}
	} else {
		if !*isParagraphOpen {
			result.WriteString(paragraphOpeningTag)
			*isParagraphOpen = true
		}
		for regex, replacement := range TagMap {
			line = regex.ReplaceAllString(line, replacement)
		}
		result.WriteString(line + "\n")
	}
}

func hasUnclosedTags(data string) bool {
	regexList := []*regexp.Regexp{
		regexp.MustCompile(`(\*\*[^*]+)`),
		regexp.MustCompile(`(_[^_]+)`),
		regexp.MustCompile(`(^|\s)` + "`\\w+`"),
	}

	for _, regex := range regexList {
		matches := regex.FindAllString(data, -1)
		for _, match := range matches {
			if match[len(match)-1:] != "`" {
				return true
			}
		}
	}

	return false
}

func ConvertMarkdownToHTML(markdown string) (string, error) {
	lines := strings.Split(markdown, "\n")
	var result strings.Builder
	isPreformattedBlock := false
	isParagraphOpen := false

	for _, line := range lines {
		if strings.HasPrefix(line, backtick) {
			processPreformattedBlock(&result, &isPreformattedBlock)
			continue
		}

		if isPreformattedBlock {
			result.WriteString(line + "\n")
		} else {
			processParagraph(&result, &isParagraphOpen, line)
		}
	}

	if isParagraphOpen {
		result.WriteString(paragraphClosingTag + "\n")
	}

	if isPreformattedBlock {
		result.WriteString(preformattedBlockClosingTag)
	}

	htmlContent := result.String()

	if hasUnclosedTags(htmlContent) {
		return "", errors.New("error: invalid markdown tags are not closed")
	}

	return htmlContent, nil
}
