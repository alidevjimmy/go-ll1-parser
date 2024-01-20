package bnf

import (
	"fmt"
	"io"
	"strings"
)

type Grammar map[string][][]string

// Example
//
// <exp> : [
//
//		["a", "<exp>", <"num"],
//	 [...]
//
// ]
// ...
type Bnf struct {
	Grammar Grammar
}

func NewBnf(reader io.Reader) (*Bnf, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	dataStr := strings.TrimSpace(string(data))
	lines := strings.Split(dataStr, "\n")
	bnf := &Bnf{
		Grammar: make(Grammar),
	}
	for _, line := range lines {
		gram := strings.Split(line, "::=")
		if len(gram) != 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		nonterm, term := strings.TrimSpace(gram[0]), strings.TrimSpace(gram[1])
		prods := strings.Split(term, "|")
		bnf.Grammar[nonterm] = make([][]string, len(prods))
		for k, p := range prods {
			p = strings.TrimSpace(p)
			prod := strings.Split(p, " ")
			bnf.Grammar[nonterm][k] = append(bnf.Grammar[nonterm][k], prod...)
		}
	}
	return bnf, nil
}

func NewBnfWithGrammer(grammar Grammar) *Bnf {
	return &Bnf{
		Grammar: grammar,
	}
}

func (b *Bnf) IsTerminal(token string) bool {
	if _, ok := b.Grammar[token]; !ok {
		return true
	}
	return false
}
