package parser

import (
	lsp "go.lsp.dev/protocol"
)

type Parser interface {
	FileExtension() string
	ParseFile(path string) interface{}
	AddDefinitions(files []string)
	Methods() []string
	Diagnostics(text string, defs []ParserDefinition) []lsp.Diagnostic
	GetDefinitions() []ParserDefinition
	CompletionItem(def ParserDefinition) (lsp.CompletionItem, error)
	GetGoToDefinition(params string) []string
}

type ParserDefinition struct {
	Name  string
	Class string `yaml:"class"`
}

type PhpClass struct {
	Namespace string
	Path      string
}

// Get all structs that implements Parser interface
func InitParsers() map[string]Parser {
	return map[string]Parser{
		"service": &Service{},
	}
}
