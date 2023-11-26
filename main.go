package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var cyan = color.New(color.FgCyan).SprintFunc()
var yellow = color.New(color.FgHiYellow).SprintFunc()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if _, err := os.Stat("books"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("books", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	app := &cli.App{
		Name: "Minstrel",
		Authors: []*cli.Author{
			{
				Name:  "Chang Sau Sheong",
				Email: "sausheong@gmail.com",
			},
		},
		Copyright: "(c) 2023 Chang Sau Sheong",
		Usage:     "Using AI to create stories",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "genre",
				Aliases: []string{"g"},
				Value:   "science-fiction",
				Usage:   "the genre of the story to create",
			},
			&cli.StringFlag{
				Name:    "style",
				Aliases: []string{"s"},
				Value:   "Isaac Asimov",
				Usage:   "the author style to generate the story with",
			},
			&cli.IntFlag{
				Name:    "chapters",
				Aliases: []string{"c"},
				Value:   5,
				Usage:   "number of chapters, must be more than 3",
				Action: func(ctx *cli.Context, v int) error {
					if v <= 3 {
						return fmt.Errorf("you have tried to set %d chapters. "+
							"Each story must have at least 4 chapters", v)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "author",
				Aliases: []string{"a"},
				Value:   "Your Friendly AI Assistant",
				Usage:   "the name of the author",
			},
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Value:   "llama2:13b",
				Usage:   "large language model to use",
			},
			&cli.BoolFlag{
				Name:    "epub",
				Aliases: []string{"b"},
				Value:   false,
				Usage:   "generate ePUB file",
			},
			&cli.BoolFlag{
				Name:    "openai",
				Aliases: []string{"o"},
				Value:   false,
				Usage:   "use OpenAI",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Value:   false,
				Usage:   "print chapters to screen",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("openai") && os.Getenv("OPENAI_API_KEY") == "" {
				log.Fatalln("Please provide OPENAI_API_KEY")
			}
			if c.Bool("epub") && os.Getenv("REPLICATE_API_TOKEN") == "" {
				log.Fatalln("Please provide REPLICATE_API_TOKEN to generate cover page for epub file")
			}
			create(c.String("model"), c.Int("chapters"),
				c.String("genre"), c.String("style"), c.String("author"),
				c.Bool("verbose"), c.Bool("epub"), c.Bool("openai"))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// main function for writing the story
func create(model string, numChapters int, genre string, authorStyle string,
	author string, print bool, epub bool, openai bool) {
	fmt.Println(yellow("Minstrel is starting, creating a StoryJSON now"))
	text := []string{}
	html := []string{}

	startTime := time.Now()
	// create the story JSON from the writing prompt
	writingPrompt, _ := readFile("writing_prompt.txt")

	startPrompt := fmt.Sprintf(string(start), authorStyle, genre)
	storyJSON, _ := generate(startPrompt, writingPrompt, model, "json", openai)
	if print {
		fmt.Println(yellow("> This is StoryJSON used in the story."))
		fmt.Println()
		fmt.Println(cyan(storyJSON))
		fmt.Println()
	}
	story := Story{}
	err := json.Unmarshal([]byte(storyJSON), &story)
	if err != nil {
		log.Fatalln("Cannot unmarshal story JSON:", err)
	}
	var imageGet string
	if epub {
		imageGet = sendCreateImage(story)
	}

	// generate the first chapter
	firstChapter, _ := generate(first, storyJSON, model, "", openai)
	fmt.Println(yellow("> chapter 1"))
	if print {
		fmt.Println()
		fmt.Println(firstChapter)
		fmt.Println()
	}
	text = append(text, firstChapter)
	html = append(html, string(mdToHTML([]byte(firstChapter))))

	// generate the second chapter
	context := storyJSON + "\n\n" + firstChapter
	nextChapter, _ := generate(next, context, model, "", openai)
	fmt.Println(yellow("> chapter 2"))
	if print {
		fmt.Println()
		fmt.Println(nextChapter)
		fmt.Println()
	}
	text = append(text, nextChapter)
	html = append(html, string(mdToHTML([]byte(nextChapter))))

	// generate the next few chapters
	for i := 0; i < numChapters-2; i++ {
		context = storyJSON + "\n\n" + nextChapter
		nextChapter, _ = generate(next, context, model, "", openai)
		fmt.Println(yellow(fmt.Sprintf("> chapter %d", i+3)))
		if print {
			fmt.Println()
			fmt.Println(nextChapter)
			fmt.Println()
		}
		text = append(text, nextChapter)
		html = append(html, string(mdToHTML([]byte(nextChapter))))
	}

	// generate the final chapter
	context = storyJSON + "\n\n" + nextChapter
	lastChapter, _ := generate(last, context, model, "", openai)
	fmt.Println(yellow("> final chapter!"))
	if print {
		fmt.Println()
		fmt.Println(lastChapter)
		fmt.Println()
	}
	text = append(text, lastChapter)
	html = append(html, string(mdToHTML([]byte(lastChapter))))

	saveFile(fmt.Sprintf("books/%s.md", story.Title), strings.Join(text, "\n\n"), false)
	saveFile(fmt.Sprintf("books/%s.html", story.Title), strings.Join(html, "\n\n"), false)
	if epub {
		fmt.Println(yellow("> generating epub file"))
		coverUrl := getImage(imageGet)
		publish(story, author, html, coverUrl)
	}
	fmt.Println(yellow(fmt.Sprintf("Done in %s", time.Since(startTime))))
}

// generate the story using either Ollama or OpenAI
func generate(prompt string, context string, model string, format string, openai bool) (string, error) {
	if openai {
		return gpt(prompt, context, "gpt-3.5-turbo-16k", format, 16*1024)
	} else {
		return local(prompt, context, model, format, 4*1024) // Llama-2 context size is 4096
	}
}

// generate the story using Ollama
func local(prompt string, context string, model string, format string, modelContextSize int) (string, error) {
	// figure out the maximum tokens to generate
	csize := tokens(context)
	psize := tokens(prompt)
	maxTokens := modelContextSize - (csize + psize + 64)

	req := &CompletionRequest{
		Model:  model,
		Prompt: prompt,
		System: context,
		Options: Options{
			Temperature: 0.9,
			NumPredict:  maxTokens,
			NumCtx:      8 * 1024,
		},
		Stream: false,
	}
	if format == "json" {
		req.Format = "json"
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err in marshaling:", err)
		return "", err
	}

	r := bytes.NewReader(reqJson)
	httpResp, err := http.Post("http://localhost:11434/api/generate", "application/json", r)
	if err != nil {
		fmt.Println("err in calling ollama:", err)
		return "", err
	}
	decoder := json.NewDecoder(httpResp.Body)
	resp := &CompletionResponse{}
	decoder.Decode(&resp)
	return resp.Response, nil
}

// generate the story using OpenAI GPT
func gpt(prompt string, context string, model string, format string, modelContextSize int) (string, error) {
	csize := tokens(context)
	psize := tokens(prompt)
	maxTokens := modelContextSize - (csize + psize + 64)

	requestURL := "https://api.openai.com/v1/chat/completions"
	context = strings.ReplaceAll(context, `"`, `\"`)
	context = strings.ReplaceAll(context, "\n", " ")

	prompt = strings.ReplaceAll(prompt, `"`, `\"`)
	prompt = strings.ReplaceAll(prompt, "\n", " ")
	prompt = strings.ReplaceAll(prompt, "\t", " ")

	jsonBody := `{
	"model": "` + model + `",
	"messages": [
		{
		"role": "system",
		"content": "` + context + `"
		},
		{
		"role": "user",
		"content": "` + prompt + `"
		}
	],
	"max_tokens" : ` + strconv.Itoa(maxTokens) + `,
	"temperature" : 0.9
}`
	bodyReader := bytes.NewReader([]byte(jsonBody))
	fmt.Println(string(jsonBody))
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))

	client := http.Client{
		Timeout: 5 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return "", err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return "", err
	}

	response := OpenAIResponse{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
