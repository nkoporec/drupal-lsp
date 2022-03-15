package langserver

import (
	"fmt"
	"io/ioutil"

	"github.com/z7zmey/php-parser/php7"
)

func PhpParser(file string) *php7.Parser {
	// Read the file contents.
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print(err)
	}

	parser := php7.NewParser(src, "7.4")
	parser.Parse()

	return parser
}
