package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// concurrently get wiki pages as markdown text

// command line option to specify the output file, default to stdout if empty
var outputFile = flag.String("output", "", "path to output file (default stdout)")

// command line option to specify a list of search words
var searchWords = flag.String("search", "", "search word, use comma to separate multiple words")

// the base URL for the wikipedia page
const baseURL = "en.wikipedia.org/wiki/"

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
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
		outputStream = f
		defer f.Close()
	}

	// get the list of search words
	words := strings.Split(*searchWords, ",")
	numWords := len(words)

}

func fetchWikiPage(word string) string {
    // build the URL
    url := baseURL + word

    // fetch the page
    resp, err := http.Get(url)

    if err != nil {
        fmt.Println("Error fetching page:", err)
        return fmt.Sprintf("HTML Status Code %v while fetching Page for %v", resp.StatusCode, word)
    }

    defer resp.Body.Close()
    rc := resp.Body

    bodyBytes, err := io.ReadAll(rc)
    if err != nil {
        fmt.Println("Error reading response body:", err)
        return ""
    }

    bodyStringHTML := string(bodyBytes)
    return bodyString
}