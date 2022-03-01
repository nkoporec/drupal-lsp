package langserver

import (
	"drupal-lsp/langserver/parser"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.lsp.dev/uri"
	"gopkg.in/yaml.v2"
)

type Indexer struct {
	DocumentRoot string
	Services     []parser.ServiceDefinition
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

		if strings.Contains(path, ".services.yml") {
			serviceFile, err := ioutil.ReadFile(path)
			service := &parser.ServiceYaml{}
			if err == nil {
				err = yaml.Unmarshal(serviceFile, &service)
			}

			for name, item := range service.Services {
				i.Services = append(i.Services, parser.ServiceDefinition{
					Name:  name,
					Class: item.Class,
				})
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
	fixed := FixDocumentRootUri(s)
	v = uri.URI(fixed)

	return v.Filename()
}
