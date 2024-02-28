package main

import (
	"flag"
	"fmt"
	"github.com/SharpDevOps10/Methodologies-Lab1/markdownconverter"
	"os"
)

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

	htmlContent, err := markdownconverter.ConvertMarkdownToHTML(string(markdownContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting Markdown to HTML: %v\n", err)
		os.Exit(1)
	}
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
