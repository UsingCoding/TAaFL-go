package app

import (
	"compiler/pkg/common/grammary"
	slr "compiler/pkg/slr/common"
	"compiler/pkg/slr/common/inlinedgrammary"
	"fmt"
)

type Generator interface {
	GenerateTable(grammar inlinedgrammary.Grammar) (slr.Table, error)
}

func NewGenerator() Generator {
	return &generator{}
}

type generator struct {
}

type generateTask struct {
	table          slr.TableMap
	tableRefs      slr.TableRefs
	tableRefsStack []slr.TableRef
}

//func ()  {
//
//}

func (g *generator) GenerateTable(grammar inlinedgrammary.Grammar) (slr.Table, error) {
	table := slr.TableMap{}
	var tableRefs slr.TableRefs
	// stack with refs that needs to be proceed
	var tableRefsStack []slr.TableRef

	// adding Axiom as first element to be proceed
	tableRefs = append(tableRefs, slr.TableEntry{slr.GrammarEntry{
		Symbol:       grammar.Axiom(),
		RuleNumber:   0,
		NumberInRule: 0,
	}})
	tableRefsStack = append(tableRefsStack, 0)

	for len(tableRefsStack) != 0 {
		fmt.Println("STACK", tableRefsStack)
		fmt.Println("TABLE REFS", tableRefs)
		tableRef := tableRefsStack[len(tableRefsStack)-1]
		tableEntry := tableRefs[tableRef]

		// pop stack
		tableRefsStack = tableRefsStack[:len(tableRefsStack)-1]

		for _, grammarEntry := range tableEntry {
			rule := grammar.Rules()[grammarEntry.RuleNumber]
			if int(grammarEntry.NumberInRule) == (len(rule.RuleSymbols())) {
				continue
			}
			nextNumberInRule := grammarEntry.NumberInRule + 1

			if grammarEntry.Symbol == grammar.Axiom() {
				nextNumberInRule = 0
			}
			fmt.Println("RULE SYMBOLS", rule.RuleSymbols())
			fmt.Println("NEXT NUMBER", nextNumberInRule)
			symbol := rule.RuleSymbols()[nextNumberInRule]

			if !symbol.NonTerminal() {
				newTableRef := newTableRef(tableRefs)
				fmt.Println("NEW TABLE REF", newTableRef)

				if _, ok := table[tableRef]; !ok {
					table[tableRef] = map[grammary.Symbol]slr.TableRef{}
				}
				table[tableRef][symbol] = newTableRef

				newTableEntry := slr.TableEntry{}

				newTableEntry = append(newTableEntry, slr.GrammarEntry{
					Symbol:       symbol,
					RuleNumber:   grammarEntry.RuleNumber,
					NumberInRule: nextNumberInRule + 1,
				})
				tableRefs = append(tableRefs, newTableEntry)
				tableRefsStack = append(tableRefsStack, newTableRef)
				continue
			}

			rulesMap := grammar.FindByLeftSideSymbol(symbol)
			fmt.Println("RULES MAP SYMBOL", symbol)
			fmt.Println("RULES MAP", rulesMap)

			for ruleNumber, rule := range rulesMap {
				const firstSymbolPos = 0
				symbol = rule.RuleSymbols()[firstSymbolPos]

				if _, ok := table[tableRef]; !ok {
					table[tableRef] = map[grammary.Symbol]slr.TableRef{}
				}

				tableRef2, ok := table[tableRef][symbol]
				if !ok {
					tableRef2 = newTableRef(tableRefs)
					table[tableRef][symbol] = tableRef2
					tableRefsStack = append(tableRefsStack, tableRef2)
				}

				newTableEntry := slr.TableEntry{}

				if !symbol.NonTerminal() {
					newTableEntry = append(newTableEntry, slr.GrammarEntry{
						Symbol:       symbol,
						RuleNumber:   ruleNumber,
						NumberInRule: firstSymbolPos + 1,
					})
				}
			}
			tableRefs = append(tableRefs, newTableEntry)
			tableRefsStack = append(tableRefsStack, newTableRef)
		}
	}

	return slr.Table{
		TableMap:  table,
		TableRefs: tableRefs,
	}, nil
}

func proceedRecursiveTransitiveClosure() {

}

func newTableRef(refs slr.TableRefs) slr.TableRef {
	return slr.TableRef(len(refs))
}
