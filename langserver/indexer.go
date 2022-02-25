package langserver

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type drupalServiceKeys struct {
	Class     string `yaml:"class"`
	Arguments string `yaml:"arguments"`
}

type Indexer struct {
	documentRoot string
}

func NewIndexer(rootUri string) *Indexer {
	return &Indexer{
		documentRoot: rootUri,
	}
}

func (i *Indexer) Run() {
	i.documentRoot = fixDocumentRootUri(i.documentRoot)

	// Check if the document root exists.
	if _, err := os.Stat(i.documentRoot); os.IsNotExist(err) {
		log.Fatal("Indexer: Directory does not exist")
	}

	// Walk the document root and get all .services.yml files.
	// @todo Make it for other file types.
	services := make([]string, 0)
	filepath.Walk(i.documentRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a .services.yml file.
		if info.IsDir() {
			return nil
		}

		if strings.Contains(path, ".services.yml") {
			services = append(services, path)
		}

		return nil
	})

	// Parse the services.
	m := make(map[string]string)
	for _, serviceFilepath := range services {
		serviceFile, err := ioutil.ReadFile(serviceFilepath)
		if err == nil {
			err = yaml.Unmarshal(serviceFile, &m)
		}

		log.Println(m)
	}
}

// Remove the file:// prefix so we can access the folder.
func fixDocumentRootUri(s string) string {
	if strings.HasPrefix(s, "file://") {
		return s[len("file://"):]
	}

	// vscode sends an additonal / at the end of the uri.
	if strings.HasPrefix(s, "file:///") {
		return s[len("file:///"):]
	}

	return s
}
