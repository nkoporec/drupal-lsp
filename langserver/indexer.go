package langserver

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nkoporec/drupal-lsp/langserver/parser"

	"go.lsp.dev/uri"
)

type Indexer struct {
	DocumentRoot string
	Parsers      []parser.Parser
	PhpClasses   []parser.PhpClass
}

func NewIndexer(rootUri string) *Indexer {
	return &Indexer{
		DocumentRoot: rootUri,
	}
}

// @todo Implement caching.
func (i *Indexer) Run() {
	i.DocumentRoot = FixDocumentRootUri(i.DocumentRoot)

	// Check if the document root exists.
	if _, err := os.Stat(i.DocumentRoot); os.IsNotExist(err) {
		log.Fatal("Indexer: Directory does not exist")
	}

	// Walk the document root and get all php files.
	go func() {
		log.Println("Indexing php files ...")

		phpClasses := []parser.PhpClass{}
		filepath.Walk(i.DocumentRoot, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			// General PHP parser.

			ext := filepath.Ext(path)
			if ext != ".php" {
				return nil
			}

			className := filepath.Base(path)
			className = strings.TrimSuffix(className, ext)

			src, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println(err)
				return nil
			}

			// Split by line.
			line := 1
			temp := strings.Split(string(src), "\n")
			for _, item := range temp {
				if strings.HasPrefix(item, "namespace") {
					// Remove the namespace.
					item = strings.Replace(item, "namespace", "", 1)
					item = strings.Replace(item, " ", "", 1)
					item = strings.Replace(item, ";", "", 1)

					// Add the namespace to the list.
					phpClasses = append(phpClasses, parser.PhpClass{
						Namespace:   item + "\\" + className,
						Path:        path,
						Description: "",
					})
				}
				line++
			}

			return nil
		})

		i.PhpClasses = phpClasses
		log.Println("Indexing php files completed.")
	}()

	// Get available parsers.
	p := parser.InitParsers()
	items := make(map[string][]string)
	// Walk the document root and get all .services.yml files.
	filepath.Walk(i.DocumentRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Get available custom parsers.
		for name, item := range p {
			if strings.Contains(path, item.FileExtension()) {
				items[name] = append(items[name], path)
			}
		}

		return nil
	})

	// Parse the files.
	for name, par := range p {
		par.AddDefinitions(items[name])
		i.Parsers = append(i.Parsers, par)
	}
}

// Remove the file:// prefix so we can access the folder.
func FixDocumentRootUri(s string) string {
	if strings.HasPrefix(s, "file://") {
		return s[len("file://"):]
	}

	// vscode sends an additonal / at the end of the uri.
	if strings.HasPrefix(s, "file:///") {
		return s[len("file:///"):]
	}

	return s
}

func UriToFilename(v uri.URI) string {
	s := string(v)
	v = uri.URI(s)

	return v.Filename()
}
