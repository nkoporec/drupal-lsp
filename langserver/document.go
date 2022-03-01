package langserver

import (
	"strings"
	"unicode"

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

	if position.Character > 0 {
		offset += int(position.Character)
	}
	o := offset
	skipOpen := 1

	for o >= 0 {
		token := doc[o]
		o--
		if token == ';' || token == '}' || token == '{' {
			break
		} else if token == '(' && skipOpen <= 0 {
			break
		} else if token == ')' {
			skipOpen++
		} else if token == '(' {
			skipOpen--
		} else if token == '\n' || token == '\r' {
			continue
		}
	}

	if o+1 > len(doc) {
		return doc[o:]
	}

	methodCallLine := doc[o+2 : o+2+(offset-o-2)]
	for i := 0; i < len(methodCallLine); i++ {
		if !unicode.IsSpace(rune(methodCallLine[i])) {
			methodCallLine = methodCallLine[i:]
			break
		}
	}

	return methodCallLine
}
