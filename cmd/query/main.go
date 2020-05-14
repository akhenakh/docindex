package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
)

var (
	version = "no version from LDFLAGS"

	indexPath          = flag.String("indexPath", "index.db", "index db path")
	queryString          = flag.String("q", "", "query string")
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

	searchResults, err := index.Search(search)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(searchResults)
}
