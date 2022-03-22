package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/nkoporec/php-parser/pkg/conf"
	"github.com/nkoporec/php-parser/pkg/errors"
	"github.com/nkoporec/php-parser/pkg/parser"
	"github.com/nkoporec/php-parser/pkg/version"

	lsp "go.lsp.dev/protocol"
	"gopkg.in/yaml.v2"
)

type Service struct {
	File        *ServiceYaml
	Definitions []ServiceDefinition
}

type ServiceYaml struct {
	Services map[string]ServiceDefinition
}

type ServiceDefinition struct {
	Name  string
	Class string `yaml:"class"`
}

func (s *Service) ParseFile(path string) interface{} {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	service := &ServiceYaml{}
	if err == nil {
		err = yaml.Unmarshal(file, &service)
	}

	return service
}

func (s *Service) AddDefinitions(items []string) {
	for _, file := range items {
		item := s.ParseFile(file)
		if item == nil {
			continue
		}

		defs := item.(*ServiceYaml)
		for name, def := range defs.Services {
			s.Definitions = append(s.Definitions, ServiceDefinition{
				Name:  name,
				Class: def.Class,
			})
		}
	}
}

func (s *Service) FileExtension() string {
	return "services.yml"
}

func (s *Service) Methods() []string {
	return []string{
		"service",
		"get",
	}
}

func (s *Service) GetDefinitions() []ServiceDefinition {
	return s.Definitions
}

func (s *Service) CompletionItem(def ServiceDefinition) (lsp.CompletionItem, error) {
	return lsp.CompletionItem{
		Kind:   lsp.VariableCompletion,
		Label:  def.Name,
		Detail: fmt.Sprintf("Class %s", def.Class),
		Documentation: lsp.MarkupContent{
			Kind:  lsp.PlainText,
			Value: def.Class,
		},
	}, nil
}

func (s *Service) Diagnostics(text string) []lsp.Diagnostic {
	result := []lsp.Diagnostic{}
	src := []byte(text)

	var parserErrors []*errors.Error
	errorHandler := func(e *errors.Error) {
		parserErrors = append(parserErrors, e)
	}

	// Parse
	rootNode, err := parser.Parse(src, conf.Config{
		Version:          &version.Version{Major: 5, Minor: 6},
		ErrorHandlerFunc: errorHandler,
	})

	if err != nil {
		log.Fatal(err)
	}

	goDumper := NewServiceDumper(os.Stdout)
	rootNode.Accept(goDumper)

	diag := lsp.Diagnostic{
		Code:     2,
		Message:  "Test",
		Source:   "drupal-lsp",
		Severity: lsp.SeverityError,
		Range: lsp.Range{
			Start: lsp.Position{
				Line:      float64(1),
				Character: float64(4),
			},
			End: lsp.Position{
				Line:      float64(1),
				Character: float64(8),
			},
		},
	}

	result = append(result, diag)

	return result
}
