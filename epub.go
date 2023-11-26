package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-shiori/go-epub"
)

// create epub file
func publish(story Story, author string, htmlSections []string, coverUrl string) {
	epub, err := epub.NewEpub(story.Title)
	if err != nil {
		log.Println("Cannot create epub:", err)
		return
	}
	epub.SetAuthor(author)
	epub.SetDescription(story.Plot)
	fmt.Println(yellow("> creating cover image"))
	coverImage, _ := epub.AddImage(coverUrl, "cover.png")
	createCoverCSS(story, author)
	coverCSS, err := epub.AddCSS("assets/cover.css", "cover.css")
	if err != nil {
		fmt.Println("Cannot add CSS:", err)
	}
	epub.SetCover(coverImage, coverCSS)

	firstSectionTemplate := `<h1>%s</h1>
	<h3>%s</h3>
	<p>%s</p>
	`
	firstSection := fmt.Sprintf(firstSectionTemplate, story.Title, author, story.Plot)
	_, err = epub.AddSection(firstSection, "Front Page", "", "")
	if err != nil {
		log.Println("cannot generate front page", err)
		return
	}
	// Add a section
	for i, section := range htmlSections {
		sectionBody := html.UnescapeString(section)
		_, err = epub.AddSection(sectionBody, fmt.Sprintf("Chapter %d", i+1), "", "")
		if err != nil {
			log.Println(err)
			return
		}
	}
	// Write the epub to file
	err = epub.Write("books/" + story.Title + ".epub")
	if err != nil {
		log.Println("Cannot create epub file:", err)
	}
}

// create the cover CSS, to overlay the title over the cover image
func createCoverCSS(story Story, author string) {
	templateCss, err := readFile("assets/cover_css.template")
	if err != nil {
		log.Println("cannot read CSS template file:", err)
		return
	}
	coverCss := fmt.Sprintf(templateCss, story.Title, author)
	saveFile("assets/cover.css", coverCss, true)
}

// send create image to Replicate, calling the SDXL model to create a cover image
func sendCreateImage(story Story) string {
	sdxlreq := ReplicateSDXLRequest{
		Version: "39ed52f2a78e934b3ba6e2a89f5b1c712de7dfea535525255b1aa35c5565e08b",
		Input: Input{
			Width:             768,
			Height:            1024,
			Prompt:            story.Plot,
			Refine:            "expert_ensemble_refiner",
			Scheduler:         "K_EULER",
			LoraScale:         0.6,
			NumOutputs:        1,
			GuidanceScale:     7.5,
			ApplyWatermark:    false,
			HighNoiseFrac:     0.8,
			NegativePrompt:    "",
			PromptStrength:    0.8,
			NumInferenceSteps: 25,
		},
	}

	reqJson, err := json.Marshal(sdxlreq)
	if err != nil {
		log.Println("err in marshaling:", err)
		return ""
	}
	r := bytes.NewReader(reqJson)
	requestURL := "https://api.replicate.com/v1/predictions"
	req, err := http.NewRequest(http.MethodPost, requestURL, r)
	if err != nil {
		log.Printf("could not create Replicate request: %s\n", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", os.Getenv("REPLICATE_API_TOKEN")))
	client := http.Client{
		Timeout: 5 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("error making http Replicate request: %s\n", err)
		return ""
	}
	response := ReplicateSDXLResponse{}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("could not read Replicate response body: %s\n", err)
		return ""
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		log.Printf("cannot unmarshal Replicate response to structs: %s\n", err)
		return ""
	}
	return response.Urls.Get
}

// get the image URL, this is called after sendCreateImage
func getImage(url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("could not create Replicate request: %s\n", err)
		return ""
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", os.Getenv("REPLICATE_API_TOKEN")))

	client := http.Client{
		Timeout: 5 * time.Minute,
	}

	response := ReplicateSDXLResponse{}
	for {
		res, err := client.Do(req)
		if err != nil {
			log.Printf("error making http Replicate request: %s\n", err)
			return ""
		}
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("could not read Replicate response body: %s\n", err)
			return ""
		}
		err = json.Unmarshal(resBody, &response)
		if err != nil {
			log.Printf("cannot unmarshal Replicate response to structs: %s\n", err)
			return ""
		}
		if response.Status == "succeeded" {
			break
		}
		time.Sleep(2 * time.Second)
	}
	return response.Output[0]
}
