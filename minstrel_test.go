package main

import (
	"encoding/json"
	"fmt"
	"html"
	"testing"
	"time"
)

func TestGetImage(t *testing.T) {

	storyJSON := `{
"title": "The Fossil's Revenge",
"plot": "As a renowned paleontologist, you've always had the gift of touching fossils and reliving the last moments of the creatures they once were. However, when you accidentally touch a particularly gruesome fossil, you find yourself trapped in a never-ending loop of reliving the same moment over and over again. As you struggle to break free from the fossil's hold, you begin to realize that the creature may not be as dead as you thought...",
"characters": ["Dr. Emma Taylor", "Fossil #17"],
"locations": ["Museum of Natural History", "Desert Canyon", "Ancient Ruins"],
"author_style": "PG Wodehouse",
"genre": "fantasy comedy"
}`
	story := Story{}
	err := json.Unmarshal([]byte(storyJSON), &story)
	if err != nil {
		t.Error("Cannot unmarshal story JSON:", err)
	}

	url := sendCreateImage(story)
	fmt.Println("get url:", url)
	time.Sleep(3 * time.Second)

	results := getImage(url)
	fmt.Println("results:", results)
}

func TestPublish(t *testing.T) {
	author := "Chang Sau Sheong"
	storyJSON := `{
		"title": "The Fossil's Revenge",
		"plot": "As a renowned paleontologist, you've always had the gift of touching fossils and reliving the last moments of the creatures they once were. However, when you accidentally touch a particularly gruesome fossil, you find yourself trapped in a never-ending loop of reliving the same moment over and over again. As you struggle to break free from the fossil's hold, you begin to realize that the creature may not be as dead as you thought...",
		"characters": ["Dr. Emma Taylor", "Fossil #17"],
		"locations": ["Museum of Natural History", "Desert Canyon", "Ancient Ruins"],
		"author_style": "PG Wodehouse",
		"genre": "fantasy comedy"
		}`
	story := Story{}
	err := json.Unmarshal([]byte(storyJSON), &story)
	if err != nil {
		t.Error("Cannot unmarshal story JSON:", err)
	}

	coverUrl := "https://replicate.delivery/pbxt/QRw8JLzaYmriOJLS8fV6xg6ta4eQftDYf8A5D2n6E2wihSvHB/out-0.png"

	html1 := `<p>Of course! Here&rsquo;s the first chapter of &ldquo;The Fossil&rsquo;s Revenge&rdquo; in the style of PG Wodehouse:</p>

<h1 id="chapter-1-a-gift-and-a-curse">Chapter 1 - A Gift and a Curse</h1>

<p>Dr. Emma Taylor was no stranger to the world of paleontology. As a renowned expert in her field, she had spent countless hours studying the ancient creatures that once roamed the earth. But despite her extensive knowledge, there was one particular fossil that continued to elude her - Fossil #17.</p>

<p>Emma had been searching for this particular fossil for years, and finally, after a grueling expedition through the desert canyons, she had found it. Or so she thought. As soon as she touched the fossil, Emma felt a strange sensation wash over her - a feeling that she was no longer in control of her own body.</p>

<p>At first, Emma thought it was just her imagination playing tricks on her. But as the moments passed and she found herself reliving the same moment over and over again, she realized that something much more sinister was at play. It seemed that Fossil #17 had a strange power over her, trapping her in a never-ending loop of reliving the creature&rsquo;s final moments.</p>

<p>As Emma tried to break free from the fossil&rsquo;s hold, she began to notice strange things happening around her. The museum&rsquo;s exhibits seemed to be changing, and the other paleontologists were acting strangely. It was as if Fossil #17 was trying to communicate with her, but in a language only it understood.</p>

<p>Emma knew that she had to find a way to break the curse of the fossil and return to her normal life. But as she delved deeper into the mystery, she realized that the creature may not be as dead as she thought&hellip;</p>


<p>Of course! Here&rsquo;s the next chapter of &ldquo;The Fossil&rsquo;s Revenge&rdquo; in the style of PG Wodehouse:</p>`

	html2 := `<h1 id="chapter-2-a-plot-unfolds">Chapter 2 - A Plot Unfolds</h1>

<p>Emma was determined to uncover the secrets of Fossil #17, despite the strange occurrences that had been plaguing her since the moment she touched it. As she delved deeper into the mystery, she discovered a sinister plot unfolding around her. It seemed that the creature trapped within the fossil was not content with just tormenting her - it had bigger plans in store for the world of paleontology.</p>

<p>As Emma navigated through the winding corridors of the museum, she stumbled upon a group of shady characters huddled around Fossil #17. They were chanting incantations and offering sacrifices to the creature within, seemingly unaware of the chaos they were causing. Emma knew that she had to stop them before it was too late, but as she approached, one of the cultists spotted her and sounded the alarm.</p>

<p>Within seconds, the museum was filled with shouting and confusion, as the cultists chased after Emma, determined to protect their dark secret. Emma knew that she had to think quickly and come up with a plan to outsmart them, but as she turned to run, she tripped over a display case and tumbled to the ground.</p>

<p>As she lay there, dazed and disoriented, Emma heard a faint whisper in her ear - it was Fossil #17, communicating with her once again. The creature seemed to be urging her on, guiding her towards some unknown destination, but Emma couldn&rsquo;t quite decipher its meaning. She knew that she had to act fast before the cultists caught up with her, but she also knew that she couldn&rsquo;t trust the whispers of a cursed fossil.</p>

<p>Just as she was trying to make sense of it all, Emma felt a sudden jolt of energy and found herself transported to an ancient ruin deep in the desert. She knew that this was no ordinary place - Fossil #17 had brought her here for some sinister purpose of its own. But before she could ponder the reason behind her arrival, she heard a faint rustling in the shadows. It was the cultists, and they were closing in fast. Emma knew that she had to act quickly if she wanted to uncover the secrets of Fossil #17 and escape alive&hellip;</p>


<p>Of course! Here&rsquo;s the next chapter of &ldquo;The Fossil&rsquo;s Revenge&rdquo; in the style of PG Wodehouse:</p>`
	htmlSections := []string{html.UnescapeString(html1), html.UnescapeString(html2)}

	publish(story, author, htmlSections, coverUrl)

}

func TestCreateCSS(t *testing.T) {
	author := "Chang Sau Sheong"
	storyJSON := `{
		"title": "The Fossil's Revenge",
		"plot": "As a renowned paleontologist, you've always had the gift of touching fossils and reliving the last moments of the creatures they once were. However, when you accidentally touch a particularly gruesome fossil, you find yourself trapped in a never-ending loop of reliving the same moment over and over again. As you struggle to break free from the fossil's hold, you begin to realize that the creature may not be as dead as you thought...",
		"characters": ["Dr. Emma Taylor", "Fossil #17"],
		"locations": ["Museum of Natural History", "Desert Canyon", "Ancient Ruins"],
		"author_style": "PG Wodehouse",
		"genre": "fantasy comedy"
		}`
	story := Story{}
	err := json.Unmarshal([]byte(storyJSON), &story)
	if err != nil {
		t.Error("Cannot unmarshal story JSON:", err)
	}
	createCoverCSS(story, author)
}
