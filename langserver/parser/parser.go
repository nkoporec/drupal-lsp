package parser

import (
	lsp "go.lsp.dev/protocol"
)

type Parser interface {
	FileExtension() string
	ParseFile(path string) interface{}
	NewParser(data interface{})
	Methods() []string
	GetDefinitions() []ServiceDefinition
	CompletionItem(def ServiceDefinition) (lsp.CompletionItem, error)
}

// Get all structs that implements Parser interface
func InitParsers() map[string]Parser {
	return map[string]Parser{
		"service": &Service{},
	}
}
