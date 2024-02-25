package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func convertMarkdownToHTML(markdown string) string {
	html := markdown

	html = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(html, "<b>$1</b>")

	html = regexp.MustCompile(`_(.*?)_`).ReplaceAllString(html, "<i>$1</i>")

	preformattedBlocks := regexp.MustCompile("```([^`]+)```").FindAllStringSubmatch(html, -1)
	for _, match := range preformattedBlocks {
		preformattedText := match[1]
		preformattedHTML := "<pre>" + preformattedText + "</pre>"
		html = strings.Replace(html, match[0], preformattedHTML, 1)
	}

	html = regexp.MustCompile("`([^`]+)`").ReplaceAllString(html, "<tt>$1</tt>")

	return html
}

func main() {
	inputPath := flag.String("in", "", "Input Markdown file path")
	outputPath := flag.String("out", "", "Output HTML file path")
	flag.Parse()

	if *inputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Input file path not provided.")
		os.Exit(1)
	}

	markdownContent, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}

	htmlContent := convertMarkdownToHTML(string(markdownContent))

	if *outputPath != "" {
		err := os.WriteFile(*outputPath, []byte(htmlContent), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to output file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(htmlContent)
	}
}
