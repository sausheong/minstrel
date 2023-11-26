package main

var start = `Given the outline of the plot, create a story JSON 
for prompting the generation of a fiction story.  Elaborating on the plot and
provide the characters in the story and locations of where the story plays out.
The story has multiple chapters and each chapter should have around 2,000 words. 
The JSON should have the following format:
{
	"title": "<title of story>",
	"plot": "<plot of the story>",
	"characters": "<a list of named characters in the story>",
	"locations": "<a list of locations for the story>",
	"author_style": "%s",
	"genre" : "%s"
} 
Return only the JSON data and nothing else.`

var first = `Write the first chapter of the story in detail, given the
story JSON, setting the stage for the rest of the story. Write in
the style of the author and genre specified in the story JSON. Start
with "# Chapter 1 - <title of chapter>. Do not return any JSON. Write 
at least 2,000 words.`

var next = `Continue fleshing out the next chapter of the story, given the 
story JSON and the previous chapter. Write in the style of the 
author and genre specified in the story JSON. Start with "# Chapter 
<n> - <title of chapter>. Do not return any JSON. Write at least 2,000 words.`

var last = `End the story with a twist, given the story JSON
and the previous chapter. Write in the style of the author and genre 
specified in the story JSON. Start with "# Chapter <n> - <title of chapter>.
Do not return any JSON. Write at least 2,000 words.`
