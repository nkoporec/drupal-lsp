package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/nkoporec/drupal-lsp/php"
	"github.com/nkoporec/drupal-lsp/utils"

	lsp "go.lsp.dev/protocol"
	"gopkg.in/yaml.v2"
)

type Service struct {
	File        *ServiceYaml
	Definitions []ParserDefinition
}

type ServiceYaml struct {
	Services map[string]ParserDefinition
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
			s.Definitions = append(s.Definitions, ParserDefinition{
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
	}
}

func (s *Service) GetDefinitions() []ParserDefinition {
	return s.Definitions
}

func (s *Service) CompletionItem(def ParserDefinition) (lsp.CompletionItem, error) {
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

func (s *Service) Diagnostics(text string, defs []ParserDefinition) []lsp.Diagnostic {
	result := []lsp.Diagnostic{}
	src := []byte(text)

	defsNames := []string{}
	for _, def := range defs {
		defsNames = append(defsNames, def.Name)
	}

	// Parse the php file.
	parsedDoc, err := php.Parse(src)
	if err != nil {
		log.Fatal(err)
	}

	// Get all \Drupal::service calls.
	for _, static := range parsedDoc.StaticCalls {
		if static.Class.Name != "Drupal" {
			continue
		}

		if !utils.InSlice(s.Methods(), static.Method.Name) {
			continue
		}

		if len(static.Args) != 1 {
			continue
		}

		arg := static.Args[0]

		// Strip quotes from a string.
		argName := strings.Trim(arg.Name, "\"")
		argName = strings.Trim(arg.Name, "'")

		// If the arg is not in the list then show the error.
		if !utils.InSlice(defsNames, argName) {
			diag := lsp.Diagnostic{
				Code:     2,
				Message:  fmt.Sprintf("Undefined service '%s'", arg.Name),
				Source:   "drupal-lsp",
				Severity: lsp.SeverityError,
				Range: lsp.Range{
					Start: lsp.Position{
						Line:      float64(arg.Position.StartLine - 1),
						Character: float64(arg.Position.StartPos),
					},
					End: lsp.Position{
						Line:      float64(arg.Position.EndLine),
						Character: float64(arg.Position.EndPos),
					},
				},
			}
			result = append(result, diag)
		}
	}

	return result
}
