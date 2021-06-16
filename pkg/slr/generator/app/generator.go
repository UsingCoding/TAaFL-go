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

func (g *generator) GenerateTable(grammar inlinedgrammary.Grammar) (slr.Table, error) {
	task := generateStrategy{
		grammar: grammar,
		table:   slr.TableMap{},
	}

	return task.do()
}

type generateStrategy struct {
	grammar   inlinedgrammary.Grammar
	table     slr.TableMap
	tableRefs slr.TableRefs
	// stack with refs that needs to be proceed
	tableRefsStack []slr.TableRef
}

func (strategy *generateStrategy) do() (slr.Table, error) {
	// adding Axiom as first element to be proceed
	strategy.tableRefs = append(strategy.tableRefs, slr.TableEntry{slr.GrammarEntry{
		Symbol:       strategy.grammar.Axiom(),
		RuleNumber:   0,
		NumberInRule: 0,
	}})
	strategy.tableRefsStack = append(strategy.tableRefsStack, 0)

	for len(strategy.tableRefsStack) != 0 {

		tableRef := strategy.tableRefsStack[len(strategy.tableRefsStack)-1]
		tableEntry := strategy.tableRefs[tableRef]

		// pop stack
		strategy.tableRefsStack = strategy.tableRefsStack[:len(strategy.tableRefsStack)-1]

		for _, grammarEntry := range tableEntry {

			strategy.printState()

			rule := strategy.grammar.Rules()[grammarEntry.RuleNumber]
			fmt.Println("RULE", rule)

			nextNumberInRule := grammarEntry.NumberInRule + 1

			if grammarEntry.Symbol == strategy.grammar.Axiom() {
				nextNumberInRule = 0
			}

			if int(nextNumberInRule) == len(rule.RuleSymbols()) {
				fmt.Println("COLLAPSING BY GRAMMAR ENTRY", grammarEntry)
				continue
			}

			symbol := rule.RuleSymbols()[nextNumberInRule]

			if !symbol.NonTerminal() {
				strategy.safeWriteToTableEntryNewGrammarEntry(
					tableRef,
					symbol,
					slr.GrammarEntry{
						Symbol:       symbol,
						RuleNumber:   grammarEntry.RuleNumber,
						NumberInRule: nextNumberInRule,
					},
				)

				continue
			}

			strategy.safeWriteToTableEntryNewGrammarEntry(
				tableRef,
				symbol,
				slr.GrammarEntry{
					Symbol:       symbol,
					RuleNumber:   grammarEntry.RuleNumber,
					NumberInRule: nextNumberInRule,
				},
			)

			strategy.proceedRecursiveTransitClosure(tableRef, symbol)
		}
		strategy.printState()
	}

	return slr.Table{
		TableMap:  strategy.table,
		TableRefs: strategy.tableRefs,
	}, nil
}

func (strategy *generateStrategy) proceedRecursiveTransitClosure(tableRefKey slr.TableRef, handledNonTerminal grammary.Symbol) {
	rulesMap := strategy.grammar.FindByLeftSideSymbol(handledNonTerminal)

	for ruleNumber, rule := range rulesMap {
		const firstSymbolPos = 0
		symbol := rule.RuleSymbols()[firstSymbolPos]

		if symbol.NonTerminal() {
			strategy.proceedRecursiveTransitClosure(tableRefKey, symbol)
			return
		}

		strategy.safeWriteToTableEntryNewGrammarEntry(
			tableRefKey,
			symbol,
			slr.GrammarEntry{
				Symbol:     symbol,
				RuleNumber: ruleNumber,
				// we start rule with 0
				NumberInRule: firstSymbolPos,
			},
		)
	}
}

func (strategy *generateStrategy) safeWriteToTableEntryNewGrammarEntry(
	tableRefKey slr.TableRef,
	symbol grammary.Symbol,
	grammarEntry slr.GrammarEntry,
) {
	if _, ok := strategy.table[tableRefKey]; !ok {
		strategy.table[tableRefKey] = map[grammary.Symbol]slr.TableRef{}
	}

	tableRef, ok := strategy.table[tableRefKey][symbol]
	if !ok {
		// if table ref not exists we create new and put to stack to proceed later
		tableRef = strategy.newTableRef(strategy.tableRefs)
		strategy.table[tableRefKey][symbol] = tableRef
		strategy.tableRefs = append(strategy.tableRefs, slr.TableEntry{})
		strategy.tableRefsStack = append(strategy.tableRefsStack, tableRef)
	}

	tableEntry := strategy.tableRefs[tableRef]

	tableEntry = append(tableEntry, grammarEntry)

	strategy.tableRefs[tableRef] = tableEntry
}

func (strategy *generateStrategy) newTableRef(refs slr.TableRefs) slr.TableRef {
	return slr.TableRef(len(refs))
}

func (strategy *generateStrategy) printState() {
	fmt.Println("STATE")

	fmt.Println("STACK", strategy.tableRefsStack)
	fmt.Println("TABLE REFS", strategy.tableRefs)
	fmt.Println("TABLE", strategy.table)

	fmt.Println("END STATE")
}
