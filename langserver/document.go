package langserver

import (
	"errors"
	"strings"

	lsp "go.lsp.dev/protocol"
)

type Document struct {
	URI  string
	Text string
}

func (d *Document) GetDiagnostics(indexer *Indexer) ([]lsp.Diagnostic, error) {
	result := []lsp.Diagnostic{}

	p := indexer.Parsers
	for _, par := range p {
		diagnostic := par.Diagnostics(d.Text, par.GetDefinitions())
		result = append(result, diagnostic...)
	}

	return result, nil
}

func (d *Document) GetMethodCall(position lsp.Position) (string, error) {
	doc := string(d.Text)
	c := doc
	currentLine := 0
	line := int(position.Line)

	// Remove everything before the current line.
	// The result is the current line and everything after it.
	for currentLine < line {
		currentLine++
		lineEnd := strings.IndexRune(c, '\n')
		if lineEnd == -1 {
			break
		}
		c = c[lineEnd+1:]
	}

	// Find where the current line ends and remove
	// everything after it.
	// We only care about the current line and the cursor position.
	currentLineEnd := strings.IndexRune(c, '\n')
	c = c[:currentLineEnd]

	c = c[:int(position.Character)]
	// Method delimeter is a char that starts a method call
	// We can have it two ways.
	// ->method()
	// ::method()
	methodStartDelimiter := strings.IndexRune(c, '>')
	if methodStartDelimiter == -1 {
		methodStartDelimiter = strings.IndexRune(c, ':')
	}

	if methodStartDelimiter <= 0 {
		return "", errors.New("Method start delimiter not found")
	}

	// Get everything after the method delimiter.
	// We add one to the position to get the cursor position so
	// we remove the actual delimeter char (: or >)
	c = c[methodStartDelimiter+1:]

	// Method name is till the first (
	methodEndDelimiter := strings.IndexRune(c, '(')
	if methodEndDelimiter == -1 {
		return "", errors.New("Method end delimeter not found")
	}

	// Get the method name
	methodName := c[:methodEndDelimiter]

	// Remove extra : as it could be \Drupal::method();
	methodName = strings.Trim(methodName, ":")

	return methodName, nil
}

func (d *Document) GetMethodParams(position lsp.Position) (string, error) {
	doc := string(d.Text)
	c := doc
	currentLine := 0
	line := int(position.Line)

	// Remove everything before the current line.
	// The result is the current line and everything after it.
	for currentLine < line {
		currentLine++
		lineEnd := strings.IndexRune(c, '\n')
		if lineEnd == -1 {
			break
		}
		c = c[lineEnd+1:]
	}

	// Find where the current line ends and remove
	// everything after it.
	// We only care about the current line and the cursor position.
	currentLineEnd := strings.IndexRune(c, '\n')
	c = c[:currentLineEnd]

	paramBef := c[:int(position.Character)]
	paramAft := c[int(position.Character):]

	paramBefStartDelimiter := strings.IndexRune(paramBef, '(')
	paramBef = paramBef[paramBefStartDelimiter+2:]

	paramAftStartDelimiter := strings.IndexRune(paramAft, ')')
	paramAft = paramAft[:paramAftStartDelimiter]

	paramBef = strings.Trim(paramBef, "'")
	paramBef = strings.Trim(paramBef, "\"")

	paramAft = strings.Trim(paramAft, "'")
	paramAft = strings.Trim(paramAft, "\"")

	return paramBef + "" + paramAft, nil
}
