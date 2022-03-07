package langserver

import (
	"drupal-lsp/langserver/parser"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.lsp.dev/uri"
)

type Indexer struct {
	DocumentRoot string
	Parsers      []parser.Parser
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

	// Get available parsers.
	p := parser.InitParsers()
	// Walk the document root and get all .services.yml files.
	// @todo Make it for other file types.
	filepath.Walk(i.DocumentRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a .services.yml file.
		if info.IsDir() {
			return nil
		}

		// Get available parsers.
		for _, item := range p {
			if strings.Contains(path, item.FileExtension()) {
				// Parse the file.
				data := item.ParseFile(path)
				if err != nil {
					log.Printf("Indexer: Error parsing %s: %s", path, err)
					return nil
				}

				// Add the parsed file to the index.
				item.NewParser(data)
			}
		}

		return nil
	})

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
