package langserver

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type ServiceYaml struct {
	Services map[string]ServiceDefinition
}

type ServiceDefinition struct {
	Class string `yaml:"class"`
}

type Indexer struct {
	DocumentRoot string
	Services     []string
}

func NewIndexer(rootUri string) *Indexer {
	return &Indexer{
		DocumentRoot: rootUri,
	}
}

// @todo Implement caching.
func (i *Indexer) Run() {
	i.DocumentRoot = fixDocumentRootUri(i.DocumentRoot)

	// Check if the document root exists.
	if _, err := os.Stat(i.DocumentRoot); os.IsNotExist(err) {
		log.Fatal("Indexer: Directory does not exist")
	}

	// Walk the document root and get all .services.yml files.
	// @todo Make it for other file types.
	services := make([]string, 0)
	filepath.Walk(i.DocumentRoot, func(path string, info os.FileInfo, err error) error {
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
	serviceNames := make([]string, 0)
	for _, serviceFilepath := range services {
		serviceFile, err := ioutil.ReadFile(serviceFilepath)
		service := &ServiceYaml{}
		if err == nil {
			err = yaml.Unmarshal(serviceFile, &service)
		}

		for name, _ := range service.Services {
			serviceNames = append(serviceNames, name)
		}
	}

	i.Services = serviceNames
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
