package ll1

import (
	"fmt"
	"strings"

	"github.com/alidevjimmy/parsers/bnf"
)

type LL1Parser struct {
	grammar     bnf.Grammar
	firstSets   map[string]map[string]bool
	followSets  map[string]map[string]bool
	parseTable  map[string]map[string][]string
	startSymbol string
	epsilon     string
}

func NewLL1Parser(grammar bnf.Grammar, startSymbol string, epsilon string) *LL1Parser {
	parser := &LL1Parser{
		grammar:     grammar,
		startSymbol: startSymbol,
		epsilon:     epsilon,
		firstSets:   make(map[string]map[string]bool),
		followSets:  make(map[string]map[string]bool),
		parseTable:  make(map[string]map[string][]string),
	}

	parser.calculateSets()

	parser.buildParseTable()

	fmt.Println("First Set:")
	for nonTerm, firstSet := range parser.firstSets {
		fmt.Printf("%s: %v\n", nonTerm, firstSet)
	}

	fmt.Println("Follow Set:")
	for nonTerminal, followSet := range parser.followSets {
		fmt.Printf("%s: %v\n", nonTerminal, followSet)
	}

	fmt.Println("LL(1) Parse Table:")
	for nonTerminal, row := range parser.parseTable {
		for terminal, production := range row {
			fmt.Printf("%s, %s: %v\n", nonTerminal, terminal, production)
		}
	}

	return parser
}

func (parser *LL1Parser) calculateSets() {
	for nonTerminal := range parser.grammar {
		parser.firstSets[nonTerminal] = make(map[string]bool)
		parser.followSets[nonTerminal] = make(map[string]bool)
	}

	for nonTerminal := range parser.grammar {
		parser.calculateFirstSet(nonTerminal)
	}

	parser.followSets[parser.startSymbol]["$"] = true

	for nonTerminal := range parser.grammar {
		parser.calculateFollowSet(nonTerminal)
	}
}

func (parser *LL1Parser) calculateFirstSet(symbol string) {
	if _, exists := parser.firstSets[symbol]; !exists {
		parser.firstSets[symbol] = make(map[string]bool)
	}

	productions := parser.grammar[symbol]

	for _, production := range productions {
		firstSymbol := production[0]

		if parser.IsTerminal(firstSymbol) {
			parser.firstSets[symbol][firstSymbol] = true
		} else if !parser.IsTerminal(firstSymbol) {
			parser.calculateFirstSet(firstSymbol)

			for terminal := range parser.firstSets[firstSymbol] {
				parser.firstSets[symbol][terminal] = true
			}
		} else {
			parser.firstSets[symbol][parser.epsilon] = true
		}
	}
}

func (parser *LL1Parser) calculateFollowSet(nonTerminal string) {
	for symbol := range parser.grammar {
		for _, production := range parser.grammar[symbol] {
			for i, currentSymbol := range production {
				if currentSymbol == nonTerminal {
					if i+1 < len(production) {
						nextSymbol := production[i+1]
						parser.addFirstSetToFollowSet(nonTerminal, nextSymbol)
					} else if i+1 == len(production) && symbol != nonTerminal {
						parser.addFollowSetToFollowSet(nonTerminal, symbol)
					}
				}
			}
		}
	}
}

func (parser *LL1Parser) addFirstSetToFollowSet(nonTerminal, nextSymbol string) {
	for terminal := range parser.firstSets[nextSymbol] {
		if terminal != parser.epsilon {
			parser.followSets[nonTerminal][terminal] = true
		}
	}
	if _, exists := parser.firstSets[nextSymbol][parser.epsilon]; exists {
		parser.addFollowSetToFollowSet(nonTerminal, nextSymbol)
	}
}

func (parser *LL1Parser) addFollowSetToFollowSet(nonTerminal, symbol string) {
	for terminal := range parser.followSets[symbol] {
		parser.followSets[nonTerminal][terminal] = true
	}
}

func (parser *LL1Parser) buildParseTable() {
	for nonTerminal := range parser.grammar {
		parser.parseTable[nonTerminal] = make(map[string][]string)
	}

	for nonTerminal, productions := range parser.grammar {
		for _, production := range productions {
			firstSet := parser.calculateProductionFirstSet(production)
			for terminal := range firstSet {
				parser.parseTable[nonTerminal][terminal] = production
			}
			if _, exists := firstSet[parser.epsilon]; exists {
				for terminal := range parser.followSets[nonTerminal] {
					parser.parseTable[nonTerminal][terminal] = production
				}
			}
		}
	}
}

func (parser *LL1Parser) calculateProductionFirstSet(production []string) map[string]bool {
	firstSet := make(map[string]bool)

	for _, symbol := range production {
		if parser.IsTerminal(symbol) {
			firstSet[symbol] = true
			break
		} else if !parser.IsTerminal(symbol) {
			for terminal := range parser.firstSets[symbol] {
				firstSet[terminal] = true
			}
			if _, exists := parser.firstSets[symbol][parser.epsilon]; !exists {
				break
			}
		} else {
			firstSet[parser.epsilon] = true
			break
		}
	}

	return firstSet
}

func (parser *LL1Parser) getFirstSet(sequence []string) map[string]bool {
	firstSet := make(map[string]bool)

	for _, symbol := range sequence {
		if parser.IsTerminal(symbol) {
			firstSet[symbol] = true
			break
		} else {

			firstSet = parser.union(firstSet, parser.firstSets[symbol])
			if !parser.isNullable(symbol) {
				break
			}
		}
	}
	return firstSet
}

func (parser *LL1Parser) isNullable(symbol string) bool {
	for _, production := range parser.grammar[symbol] {
		if len(production) == 0 || (len(production) == 1 && production[0] == parser.epsilon) {
			return true
		}
	}

	return false
}

func (parser *LL1Parser) union(set1, set2 map[string]bool) map[string]bool {
	result := make(map[string]bool)

	for key := range set1 {
		result[key] = true
	}

	for key := range set2 {
		result[key] = true
	}

	return result
}

func (parser *LL1Parser) Parse(tokens string) ([]string, error) {
	stack := []string{"$", parser.startSymbol}
	output := make([]string, 0)
	currentIndex := 0
	tokens += " $"
	tokensSlice := strings.Split(tokens, " ")

	for len(stack) > 0 {
		top := stack[len(stack)-1]

		input := tokensSlice[currentIndex]
		fmt.Printf("Stack: %v, Input: %s\n", stack, input)

		if parser.IsTerminal(top) {
			if top == input {
				stack = stack[:len(stack)-1]
				output = append(output, input)
				if currentIndex < len(tokensSlice)-1 {
					currentIndex += 1
				}
				input = tokensSlice[currentIndex]
			} else {
				return nil, fmt.Errorf("mismatched terminal: expected %s, got %s", top, input)
			}
		} else if !parser.IsTerminal(top) {
			if production, exists := parser.parseTable[top][input]; exists {
				stack = stack[:len(stack)-1]
				if production[0] != "ε" {
					stack = append(stack, reverseProduction(production)...)
				}
			} else {
				return nil, fmt.Errorf("no production for non-terminal %s and input %s", top, input)
			}
		} else {
			return nil, fmt.Errorf("invalid symbol on the stack: %s", top)
		}
	}

	return output, nil
}

func reverseProduction(production []string) []string {
	// Reverse the production for proper stack handling
	reversed := make([]string, len(production))
	for i, symbol := range production {
		reversed[len(production)-1-i] = symbol
	}
	return reversed
}

func (parser *LL1Parser) IsTerminal(token string) bool {
	if _, ok := parser.grammar[token]; !ok {
		return true
	}
	return false
}
