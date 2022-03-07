package parser

import (
	"fmt"
	"io/ioutil"

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

func (s *Service) NewParser(data interface{}) {
	items := data.(*ServiceYaml)

	definition := &ServiceDefinition{}
	for name, item := range items.Services {
		definition.Name = name
		definition.Class = item.Class
		s.Definitions = append(s.Definitions, *definition)
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

// @todo: Old, remove once refactoring is done
func CompletionItemForService(s ServiceDefinition) (lsp.CompletionItem, error) {
	return lsp.CompletionItem{
		Kind:   lsp.VariableCompletion,
		Label:  s.Name,
		Detail: fmt.Sprintf("Class %s", s.Class),
		Documentation: lsp.MarkupContent{
			Kind:  lsp.PlainText,
			Value: s.Class,
		},
	}, nil
}
