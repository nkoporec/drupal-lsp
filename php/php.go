package php

import (
	"io"
	"log"
	"os"

	"github.com/z7zmey/php-parser/pkg/ast"
	"github.com/z7zmey/php-parser/pkg/conf"
	"github.com/z7zmey/php-parser/pkg/parser"
	"github.com/z7zmey/php-parser/pkg/version"
)

type PhpExpression struct {
	writer        io.Writer
	indent        int
	withTokens    bool
	withPositions bool
}

func Parse(src []byte) (*PhpDumper, error) {
	rootNode, err := parser.Parse(src, conf.Config{
		Version: &version.Version{Major: 5, Minor: 6},
	})

	if err != nil {
		return nil, err
	}

	phpDumper := NewPhpDumper(os.Stdout)
	rootNode.Accept(phpDumper)

	for _, expr := range phpDumper.Expressions {
		switch expr.Class.(type) {
		case *ast.NameFullyQualified:
			class := expr.Class.(*ast.NameFullyQualified)
			classParts := class.Parts
			for _, part := range classParts {
				classPart := part.(*ast.NamePart)
				log.Println("Class called: ", string(classPart.Value))
			}

		case *ast.Name:
			class := expr.Class.(*ast.Name)
			classParts := class.Parts
			for _, part := range classParts {
				classPart := part.(*ast.NamePart)
				log.Println("Class called: ", string(classPart.Value))
			}
		}

		switch expr.Call.(type) {
		case *ast.Identifier:
			methodName := string(expr.Call.(*ast.Identifier).Value)
			log.Println("Method name: ", methodName)
		}

		// Argument
		for _, item := range expr.Args {
			switch item.(type) {
			case *ast.Argument:
			}
		}
	}

	return phpDumper, nil
}
