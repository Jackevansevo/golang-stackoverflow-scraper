package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

type question struct {
	lang, text, url string
}

var url string = "http://stackoverflow.com/questions/tagged/"

var wg sync.WaitGroup

func main() {

	// [TODO] Scrape new questions every minute

	// Define some colors
	colors := [...]color.Attribute{
		color.FgYellow,
		color.FgRed,
		color.FgGreen,
		color.FgBlue,
		color.FgMagenta,
		color.FgCyan,
	}

	tags := os.Args[1:]

	// Create a random seed
	rand.Seed(time.Now().UTC().UnixNano())

	// Shuffle our tags
	for i := range tags {
		j := rand.Intn(i + 1)
		tags[i], tags[j] = tags[j], tags[i]
	}

	// Map each programming language to a color
	lang_colors := make(map[string]color.Attribute)

	for tag := range tags {
		// Check if index is out of bounds
		if tag >= len(colors) {
			// Assign a random colour
			lang_colors[tags[tag]] = colors[rand.Intn(len(colors))]
		} else {
			lang_colors[tags[tag]] = colors[tag]
		}
	}

	wg.Add(len(tags))

	// Make a question Channel
	questions := make(chan question)

	// Keep track of visited questions
	var visited = make(map[question]bool)

	go func() {
		wg.Wait()
		close(questions)
	}()

	for index := range tags {
		lang := tags[index]
		go getQuestions(lang, questions)
	}

	// Update out track of questions
	for i := range questions {
		_ = i
		question := <-questions
		if !visited[question] {
			visited[question] = true
		} else {
			fmt.Println("already visited")
		}
	}

	for key, _ := range visited {
		// Set random tag color
		c := color.New(lang_colors[key.lang])
		c.Printf("%s: ", key.lang)

		// Disable color printing for the actual Question Text
		c.DisableColor()
		fmt.Printf("%s ", key.text)

		// Print the URL
		c = color.New(color.FgCyan).Add(color.Underline)
		c.Println(key.url)
	}
	fmt.Println("Done")

}

func getQuestions(lang string, questions chan question) {
	doc := scrapePage(url + lang)
	doc.Find(".summary h3 .question-hyperlink").Each(func(i int, s *goquery.Selection) {
		r, _ := regexp.Compile("^/questions/\\d+")
		href, _ := s.Attr("href")
		url = "stackoverflow.com" + r.FindString(href)
		text := s.Text()
		questions <- question{lang, text, url}
	})
	wg.Done()
}

func scrapePage(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}
