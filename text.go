package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/pkoukk/tiktoken-go"
	tiktoken_loader "github.com/pkoukk/tiktoken-go-loader"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

// read a file and return the contents as a string
func readFile(filepath string) (string, error) {
	ba, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalln("Cannot read file:", err)
		return "", err
	}
	return string(ba), nil
}

func saveFile(filepath string, content string, overwrite bool) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Cannot save file:", err)
		return err
	}
	defer file.Close()
	if overwrite {
		file.Truncate(0)
	}
	file.WriteString(content)
	return nil
}

// get the number of tokens
func tokens(text string) int {
	tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())
	tke, _ := tiktoken.GetEncoding("cl100k_base")
	token := tke.Encode(text, nil, nil)
	return len(token)
}

// split the text given the separator, each segment must be max number of tokens
func split(text string, separator string, maxTokens int) ([]string, error) {
	parts := strings.Split(text, separator)
	sizes := make(map[string]int)
	for _, part := range parts {
		length := tokens(part)
		sizes[part[:20]] = length
	}

	for _, v := range sizes {
		if v > maxTokens {
			err := fmt.Errorf("text chunk is more than max tokens: %v", sizes)
			return []string{}, err
		}
	}

	return parts, nil
}

// chunk the text by paragraphs, given the size of the chunk and the buffer size
func chunk(text string, size int, buffer int) []string {
	chunks := strings.Split(text, "\n\n")
	grouped := ""
	groupedChunks := []string{}
	for _, chunk := range chunks {
		if len(grouped)+len(chunk) < (size - buffer) {
			grouped = grouped + "\n" + chunk
		} else {
			groupedChunks = append(groupedChunks, grouped)
			grouped = ""
		}
	}
	groupedChunks = append(groupedChunks, grouped)
	return groupedChunks
}
