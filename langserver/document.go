package langserver

import (
	"log"
	"strings"

	lsp "go.lsp.dev/protocol"
)

type Document struct {
	URI  string
	Text string
}

func (d *Document) GetMethodCall(position lsp.Position) string {
	doc := string(d.Text)
	c := doc
	currentLine := 0
	offset := 0
	line := int(position.Line)

	// @todo: this is copied over, not sure if we need to do this.
	for currentLine < line && offset < len(doc) {
		currentLine++
		lineEnd := strings.IndexRune(c, '\n')
		if lineEnd == -1 {
			break
		}
		offset += lineEnd
		if len(c) < lineEnd+1 {
			break
		}
		offset++
		c = c[lineEnd+1:]
	}

	currentLineEnd := strings.IndexRune(c, '\n')
	c = c[:currentLineEnd]

	if len(c) < int(position.Character) {
		return ""
	}

	// ->method()
	// ::method()
	methodDelimiter := strings.IndexRune(c, '>')
	if methodDelimiter == -1 {
		methodDelimiter = strings.IndexRune(c, ':')
	}

	if methodDelimiter == -1 {
		return ""
	}

	if methodDelimiter+1 > len(c) {
		return ""
	}

	c = c[methodDelimiter+1:]

	// Method name is till the first (
	methodNameEnd := strings.IndexRune(c, '(')
	methodName := c[:methodNameEnd]

	// Remove extra :
	if strings.HasPrefix(methodName, ":") {
		methodName = methodName[1:]
	}

	log.Println("Method name:", methodName)

	return methodName
}
