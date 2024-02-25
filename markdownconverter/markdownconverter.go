package markdownconverter

import (
	"regexp"
	"strings"
)

var TagMap = map[*regexp.Regexp]string{
	regexp.MustCompile(`\*\*(.*?)\*\*`): "<b>$1</b>",
	regexp.MustCompile(`_(.*?)_`):       "<i>$1</i>",
	regexp.MustCompile("`([^`]+)`"):     "<tt>$1</tt>",
}

var (
	preformattedBlockOpeningTag = "<pre>\n"
	preformattedBlockClosingTag = "</pre>\n"
	backtick                    = "```"
)

func processPreformattedBlock(result *strings.Builder, isPreformattedBlock *bool) {
	*isPreformattedBlock = !*isPreformattedBlock
	if *isPreformattedBlock {
		result.WriteString(preformattedBlockOpeningTag)
	} else {
		result.WriteString(preformattedBlockClosingTag)
	}
}

func processRegularLine(result *strings.Builder, isPreformattedBlock bool, line string) {
	if isPreformattedBlock {
		result.WriteString(line + "\n")
	} else {
		for regex, replacement := range TagMap {
			line = regex.ReplaceAllString(line, replacement)
		}
		result.WriteString(line + "\n")
	}
}

func ConvertMarkdownToHTML(markdown string) string {
	lines := strings.Split(markdown, "\n")
	var result strings.Builder
	isPreformattedBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, backtick) {
			processPreformattedBlock(&result, &isPreformattedBlock)
			continue
		}
		processRegularLine(&result, isPreformattedBlock, line)
	}

	if isPreformattedBlock {
		result.WriteString(preformattedBlockClosingTag)
	}

	return result.String()
}
