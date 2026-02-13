package parser

import (
	"strings"

	"mvdan.cc/sh/v3/syntax"
)

func Parse(input string) (*syntax.File, error) {
	p := syntax.NewParser()
	f, err := p.Parse(strings.NewReader(input), "")
	if err != nil {
		return nil, err
	}
	return f, nil
}
