package parser

import (
	lsp "go.lsp.dev/protocol"
)

type Parser interface {
	FileExtension() string
	ParseFile(path string) interface{}
	AddDefinitions(files []string)
	Methods() []string
	Diagnostics(text string) []lsp.Diagnostic
	GetDefinitions() []ServiceDefinition
	CompletionItem(def ServiceDefinition) (lsp.CompletionItem, error)
}

// Get all structs that implements Parser interface
func InitParsers() map[string]Parser {
	return map[string]Parser{
		"service": &Service{},
	}
}
