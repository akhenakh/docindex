package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/analysis/analyzer/simple"
	_ "github.com/blevesearch/bleve/analysis/analyzer/standard"
	"github.com/pkg/browser"
)

const dvURL = "https://docs.dv.nyt.net/"

var (
	version = "no version from LDFLAGS"

	indexPath   = flag.String("indexPath", "index.db", "index db path")
	queryString = flag.String("q", "", "query string")
	lucky       = flag.Bool("lucky", false, "I feel lucky, open the first found result")
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
	if *lucky {
		search.Size = 1
	}
	searchResults, err := index.Search(search)
	if err != nil {
		log.Fatal(err)
	}

	if *lucky && searchResults.Total > 0 {
		fmt.Println(searchResults.Hits[0].ID)
		path := strings.TrimSuffix(searchResults.Hits[0].ID, ".md")
		url := dvURL + path + "/"
		err = browser.OpenURL(url)
		if err != nil {
			log.Printf("can't open URL: %s error: %s", url, err)
		}
	}

	fmt.Println(searchResults)
}
