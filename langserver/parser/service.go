package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/z7zmey/php-parser/pkg/conf"
	"github.com/z7zmey/php-parser/pkg/errors"
	"github.com/z7zmey/php-parser/pkg/parser"
	"github.com/z7zmey/php-parser/pkg/version"

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

type ServiceError struct {
	Name      string
	StartLine int
	StartPos  int
	EndLine   int
	EndPos    int
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

func (s *Service) Name() string {
	return "service"
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

func (s *Service) Diagnostics(text string, defs []ServiceDefinition) []lsp.Diagnostic {
	result := []lsp.Diagnostic{}
	src := []byte(text)

	defsNames := []string{}
	for _, def := range defs {
		defsNames = append(defsNames, def.Name)
	}

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

	for _, item := range goDumper.StaticCalls {
		methodAndArg := getMethodAndArg(src, item.StartLine, item.EndLine, item.StartPos, item.EndPos)
		for method, arg := range methodAndArg {
			if inSlice(s.Methods(), method) {
				if !inSlice(defsNames, arg) {
					diag := lsp.Diagnostic{
						Code:     2,
						Message:  "Service not found",
						Source:   "drupal-lsp",
						Severity: lsp.SeverityError,
						Range: lsp.Range{
							Start: lsp.Position{
								Line: float64(item.StartLine - 1),
							},
							End: lsp.Position{
								Line: float64(item.EndLine),
							},
						},
					}
					result = append(result, diag)
				}
			}
		}
	}

	return result
}

func getMethodAndArg(text []byte, startLine int, endLine int, startPos int, endPos int) map[string]string {
	result := make(map[string]string)
	doc := string(text)
	c := doc[startPos:endPos]

	methodStartDelimiter := strings.IndexRune(c, ':')
	c = c[methodStartDelimiter+2:]

	methodEndDelimiter := strings.IndexRune(c, '(')

	// Get the method name
	methodName := c[:methodEndDelimiter]

	// Get the method arguments
	arg := c[methodEndDelimiter+2:]
	parts := strings.Split(arg, ")")
	arg = parts[0]

	// Remove the quotes
	arg = strings.Replace(arg, "'", "", -1)
	arg = strings.Replace(arg, "\"", "", -1)

	result[methodName] = arg

	return result
}

// @todo Move it to the appropriate place.
func inSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}
