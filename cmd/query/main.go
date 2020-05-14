package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/analysis/analyzer/keyword"
	_ "github.com/blevesearch/bleve/analysis/analyzer/standard"
	"github.com/pkg/browser"
)

const dvURL = "https://docs.dv.nyt.net/"

var (
	version = "no version from LDFLAGS"

	indexPath   = flag.String("indexPath", "index.db", "index db path")
	queryString = flag.String("q", "", "query string")
	lucky       = flag.Bool("lucky", false, "open the first found link")
	limit       = flag.Int("limit", 10, "limit response count")
)

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	index, err := bleve.Open(*indexPath)
	if err != nil {
		log.Fatal(err)
	}
	query := bleve.NewQueryStringQuery(*queryString)
	search := bleve.NewSearchRequest(query)
	//search.Highlight = bleve.NewHighlightWithStyle(ansi.Name)
	search.Highlight = bleve.NewHighlight()
	search.Size = *limit

	searchResults, err := index.Search(search)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(searchResults)

	if *lucky && searchResults.Total > 0 {
		fmt.Println(searchResults.Hits[0].ID)
		path := strings.TrimSuffix(searchResults.Hits[0].ID, ".md")
		log.Println("opening", dvURL+path+"/")
		browser.OpenURL(dvURL + path + "/")
	}
}
