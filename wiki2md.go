package main

import (
	"flag"
	"fmt"
	"github.com/JohannesKaufmann/html-to-markdown"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// concurrently get wiki pages as markdown text

// command line option to specify the output file, default to stdout if empty
var outputFile = flag.String("output", "", "path to output file (default stdout)")

// command line option to specify a list of search words
var searchWords = flag.String("search", "", "search word, use comma to separate multiple words")

// the base URL for the wikipedia page
const wikiAPIBaseURL = "https://en.wikipedia.org/w/api.php"

func main() {
	// parse the command line options
	flag.Parse()

	var outputStream io.Writer
	ch := make(chan string)

	if *outputFile == "" {
		outputStream = os.Stdout
	} else {

		f, err := os.OpenFile(*outputFile, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			fmt.Println("Error opening file:", *outputFile)
			log.Fatal(err)
		}
		outputStream = f
		defer f.Close()
	}

	// get the list of search words
	wordStr := strings.ReplaceAll(*searchWords, " ", "")
	words := strings.Split(wordStr, ",")
	numWords := len(words)

	// fetch the wiki pages concurrently
	for _, word := range words {
		go fetchAndConvertWikiPage(word, ch)
	}

	// write the markdown to the output file
	for i := 0; i < numWords; i++ {
		fmt.Fprintln(outputStream, <-ch)
	}

}

func buildSearchURL(word string) string {
    params := url.Values{}
    params.Add("action", "parse")
    params.Add("format", "json")
    params.Add("page", word)
    params.Add("prop", "text")
    params.Add("formatversion", "2")

    return fmt.Sprintf("%s?%s", wikiAPIBaseURL, params.Encode())
}

func fetchWikiPage(word string) string {
	// build the URL
	url := buildSearchURL(word)

	// fetch the page
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error fetching page:", url)
		log.Fatal(err)
	}

	defer resp.Body.Close()
	rc := resp.Body

	bodyBytes, err := io.ReadAll(rc)
	if err != nil {
		fmt.Println("Error reading response body")
		return ""
	}

	return string(bodyBytes)
}

func convertWikiPageToMarkdown(page string) string {
	// convert the page to markdown
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(page)

	if err != nil {
		fmt.Println("Error converting page to markdown")
		return ""
	}

	return markdown
}

func fetchAndConvertWikiPage(word string, ch chan string) {
	page := fetchWikiPage(word)
	markdown := convertWikiPageToMarkdown(page)
	ch <- markdown
}
