package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/analysis/analyzer/standard"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

var (
	version   = "no version from LDFLAGS"
	indexPath = flag.String("indexPath", "index.db", "index db path")
	docPath   = flag.String("docPath", "mkdocs.yml", "mkdocs config path")
)

type MKDocs struct {
	Nav []map[string]interface{} `yaml:"nav"`
}

type Page struct {
	Title string
	Tags  string
	Path  string
	Body  string
	itags []string
}

type txtRenderer struct {
	Title   string
	inTitle bool
}

// RenderNode is the main rendering method. It will be called once for
// every leaf node and twice for every non-leaf node (first with
// entering=true, then with entering=false). The method should write its
// rendition of the node to the supplied writer w.
func (r *txtRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.Heading:
		if node.Level == 1 {
			r.inTitle = entering
		}
	case blackfriday.Text, blackfriday.CodeBlock, blackfriday.Code:
		if node.Parent.Type == blackfriday.Link {
			// do not index self links ?
			break
		}
		if r.inTitle {
			if r.Title == "" {
				r.Title = string(node.Literal)
			} else {
				r.Title = r.Title + " " + string(node.Literal)
			}

			break
		}
		w.Write(node.Literal)
		w.Write([]byte{' '})
	}
	return blackfriday.GoToNext
}

// RenderHeader is a method that allows the renderer to produce some
// content preceding the main body of the output document. The header is
// understood in the broad sense here. For example, the default HTML
// renderer will write not only the HTML document preamble, but also the
// table of contents if it was requested.
//
// The method will be passed an entire document tree, in case a particular
// implementation needs to inspect it to produce output.
//
// The output should be written to the supplied writer w. If your
// implementation has no header to write, supply an empty implementation.
func (r *txtRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {}

// RenderFooter is a symmetric counterpart of RenderHeader.
func (r *txtRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {}

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	nav, err := readMKDocs(*docPath)
	if err != nil {
		log.Fatal(err)
	}

	var pages []Page
	for _, v := range nav.Nav {
		parseNode(v, nil, &pages)
	}

	mapping := bleve.NewIndexMapping()
	mapping.DefaultType = "Doc"

	docMapping := bleve.NewDocumentMapping()
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = standard.Name
	docMapping.AddFieldMappingsAt("Title", textFieldMapping)
	docMapping.AddFieldMappingsAt("Body", textFieldMapping)

	tagFieldMapping := bleve.NewTextFieldMapping()
	tagFieldMapping.Analyzer = simple.Name
	docMapping.AddFieldMappingsAt("Tags", tagFieldMapping)

	pathMapping := bleve.NewDocumentDisabledMapping()
	docMapping.AddSubDocumentMapping("Path", pathMapping)
	mapping.AddDocumentMapping("Doc", docMapping)

	index, err := bleve.New(*indexPath, mapping)
	if err != nil {
		log.Fatal(err)
	}

	for _, page := range pages {
		// remote link
		if strings.HasPrefix(page.Path, "http://") || strings.HasPrefix(page.Path, "https://") {
			continue
		}

		body, err := ioutil.ReadFile(path.Dir(*docPath) + "/docs/" + page.Path)
		if err != nil {
			log.Fatal(err)
		}
		renderer := &txtRenderer{}

		txt := blackfriday.Run(body, blackfriday.WithRenderer(renderer))

		page.Body = string(txt)
		page.Title = renderer.Title
		if page.Title == "" {
			page.Title = strings.Join(page.itags, " ")
		}
		page.Tags = strings.Join(page.itags, " ")
		log.Printf("Indexing: [%s] %s [%s]\n", page.Path, page.Title, strings.Join(page.itags, ","))

		err = index.Index(page.Path, page)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func parseNode(n interface{}, ptitles []string, pages *[]Page) {
	switch t := n.(type) {
	case []interface{}:
		for _, v := range t {
			parseNode(v, ptitles, pages)
		}

	case map[interface{}]interface{}:
		for title, v := range t {
			if titlestring, ok := title.(string); ok {
				ctitles := append([]string(nil), ptitles...)
				ctitles = append(ctitles, titlestring)
				parseNode(v, ctitles, pages)
				continue
			}

			parseNode(v, ptitles, pages)
		}
	case map[string]interface{}:
		for title, v := range t {
			ctitles := append([]string(nil), ptitles...)
			ctitles = append(ctitles, title)
			parseNode(v, ctitles, pages)
		}

	case string:
		if len(t) < 2 {
			return
		}
		//log.Println("filename", t, "itags", strings.Join(ptitles, "|"))
		p := Page{
			itags: ptitles,
			Path:  t,
		}
		*pages = append(*pages, p)

	default:
		log.Printf("unkown type %T %s\n", t, t)
	}
}

func readMKDocs(path string) (*MKDocs, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open mkdocs yaml file: %w", err)
	}
	doc := &MKDocs{}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("can't open mkdocs yaml file: %w", err)
	}

	err = yaml.Unmarshal(b, &doc)
	if err != nil {
		return nil, fmt.Errorf("can't read the mkdocs yaml file: %w", err)
	}
	return doc, nil
}
