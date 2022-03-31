package php

import (
	"os"

	"github.com/z7zmey/php-parser/pkg/ast"
	"github.com/z7zmey/php-parser/pkg/conf"
	"github.com/z7zmey/php-parser/pkg/parser"
	"github.com/z7zmey/php-parser/pkg/position"
	"github.com/z7zmey/php-parser/pkg/version"
)

type PhpClass struct {
	Position *position.Position
	Name     string
}

type PhpClassMethod struct {
	Position *position.Position
	Name     string
}

type PhpClassArgument struct {
	Position *position.Position
	Name     string
}

type PhpStaticCall struct {
	Position *position.Position
	Class    *PhpClass
	Method   *PhpClassMethod
	Args     []*PhpClassArgument
}

type ParsedDoc struct {
	StaticCalls []*PhpStaticCall
}

func Parse(src []byte) (*ParsedDoc, error) {
	rootNode, err := parser.Parse(src, conf.Config{
		Version: &version.Version{Major: 5, Minor: 6},
	})

	if err != nil {
		return nil, err
	}

	phpDumper := NewPhpDumper(os.Stdout)
	rootNode.Accept(phpDumper)

	// Create new parsed doc
	parsedDoc := &ParsedDoc{}

	for _, expr := range phpDumper.Expressions {
		staticCall := &PhpStaticCall{}
		isStaticCall := false

		switch expr.Class.(type) {
		// \Drupal::service('foo')
		case *ast.NameFullyQualified:
			class := expr.Class.(*ast.NameFullyQualified)
			className := ""
			classParts := class.Parts
			for _, part := range classParts {
				classPart := part.(*ast.NamePart)
				className = string(classPart.Value)
			}

			isStaticCall = true
			staticCall.Class = &PhpClass{
				Position: expr.Class.GetPosition(),
				Name:     className,
			}
		case *ast.Name:
			class := expr.Class.(*ast.Name)
			className := ""
			classParts := class.Parts
			for _, part := range classParts {
				classPart := part.(*ast.NamePart)
				className = string(classPart.Value)
			}

			isStaticCall = true
			staticCall.Class = &PhpClass{
				Position: expr.Class.GetPosition(),
				Name:     className,
			}
		}

		if isStaticCall {
			switch expr.Call.(type) {
			case *ast.Identifier:
				staticCall.Method = &PhpClassMethod{
					Position: expr.Call.GetPosition(),
					Name:     string(expr.Call.(*ast.Identifier).Value),
				}
			}

			// Argument
			for _, item := range expr.Args {
				switch item.(type) {
				case *ast.Argument:
					arg := item.(*ast.Argument)
					switch arg.Expr.(type) {
					case *ast.ScalarString:
						ex := arg.Expr.(*ast.ScalarString)
						staticCall.Args = append(staticCall.Args, &PhpClassArgument{
							Position: arg.Expr.GetPosition(),
							Name:     string(ex.Value),
						})
					}
				}
			}

			parsedDoc.StaticCalls = append(parsedDoc.StaticCalls, staticCall)
		}

	}

	return parsedDoc, nil
}
