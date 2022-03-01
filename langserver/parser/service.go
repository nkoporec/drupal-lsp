package parser

import (
	"fmt"

	lsp "go.lsp.dev/protocol"
)

type ServiceYaml struct {
	Services map[string]ServiceDefinition
}

type ServiceDefinition struct {
	Name  string
	Class string `yaml:"class"`
}

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
