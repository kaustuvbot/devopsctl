package terraform

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Parser struct {
	parser *hclparse.Parser
}

func NewParser() *Parser {
	return &Parser{
		parser: hclparse.NewParser(),
	}
}

func (p *Parser) ParseFile(filename string) (*hcl.File, error) {
	file, diags := p.parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		return nil, diags
	}
	return file, nil
}