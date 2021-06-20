package app

import (
	"compiler/pkg/common/grammary"
	slr "compiler/pkg/slr/common"
	"compiler/pkg/slr/common/inlinedgrammary"
	"fmt"
	"sort"
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
	// counter to resolve new TableRef
	tableRefsCounter *slr.TableRef
	// store proceeded table refs
	nonValidTableRefs []slr.TableRef
}

func (strategy *generateStrategy) do() (slr.Table, error) {
	// adding Axiom as first element to be proceed
	strategy.tableRefs = append(strategy.tableRefs, slr.TableEntry{
		GrammarEntries: []slr.GrammarEntry{{
			Symbol:       strategy.grammar.Axiom(),
			RuleNumber:   0,
			NumberInRule: 0,
		}},
	})
	strategy.tableRefsStack = append(strategy.tableRefsStack, strategy.newTableRef())

	for len(strategy.tableRefsStack) != 0 {

		tableRef := strategy.tableRefsStack[len(strategy.tableRefsStack)-1]
		tableEntry := strategy.tableRefs[tableRef]

		// pop stack
		strategy.tableRefsStack = strategy.tableRefsStack[:len(strategy.tableRefsStack)-1]

		for _, grammarEntry := range tableEntry.GrammarEntries {

			rule := strategy.grammar.Rules()[grammarEntry.RuleNumber]

			nextNumberInRule := grammarEntry.NumberInRule + 1

			if grammarEntry.Symbol == strategy.grammar.Axiom() {
				nextNumberInRule = 0
			}

			if int(nextNumberInRule) == len(rule.RuleSymbols()) {
				strategy.recursivelyFindCollapsingEntry(tableRef, grammarEntry)
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
	}

	return slr.Table{
		TableMap:          strategy.table,
		TableRefs:         strategy.tableRefs,
		NonValidTableRefs: strategy.nonValidTableRefs,
	}, nil
}

func (strategy *generateStrategy) proceedRecursiveTransitClosure(tableRefKey slr.TableRef, handledNonTerminal grammary.Symbol) {
	rulesMap := strategy.grammar.FindByLeftSideSymbol(handledNonTerminal)

	keys := make([]int, 0, len(rulesMap))
	for k := range rulesMap {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	for _, key := range keys {
		ruleNumber := uint(key)
		rule := rulesMap[ruleNumber]
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

func (strategy *generateStrategy) recursivelyFindCollapsingEntry(tableRefKey slr.TableRef, grammarEntry slr.GrammarEntry) {
	collapsingRule := strategy.grammar.Rules()[grammarEntry.RuleNumber]
	symbolCollapsingTo := collapsingRule.LeftSideSymbol()
	countOfSymbolsInCollapsingRule := uint(len(collapsingRule.RuleSymbols()))

	if collapsingRule.LeftSideSymbol() == strategy.grammar.Axiom() && grammarEntry.NumberInRule+1 == countOfSymbolsInCollapsingRule {
		strategy.safeWriteToTableEntryNewCollapseEntry(
			tableRefKey,
			grammary.NewSymbol(grammary.EndOfSequence),
			slr.CollapseEntry{
				RuleNumber:           grammarEntry.RuleNumber,
				Symbol:               symbolCollapsingTo,
				CountOfSymbolsInRule: countOfSymbolsInCollapsingRule,
			},
		)
		return
	}

	grammarEntries := strategy.findGrammarEntriesForSymbol(symbolCollapsingTo)

	for _, entry := range grammarEntries {
		rule := strategy.grammar.Rules()[entry.RuleNumber]
		countOfSymbolsInRule := uint(len(rule.RuleSymbols()))

		// checking that element last in axiom rule
		if rule.LeftSideSymbol() == strategy.grammar.Axiom() && entry.NumberInRule+1 == countOfSymbolsInRule {
			strategy.safeWriteToTableEntryNewCollapseEntry(
				tableRefKey,
				grammary.NewSymbol(grammary.EndOfSequence),
				slr.CollapseEntry{
					RuleNumber:           grammarEntry.RuleNumber,
					Symbol:               rule.LeftSideSymbol(),
					CountOfSymbolsInRule: countOfSymbolsInRule,
				},
			)
			continue
		}

		nextSymbol := rule.RuleSymbols()[entry.NumberInRule+1]

		if !nextSymbol.NonTerminal() {
			strategy.safeWriteToTableEntryNewCollapseEntry(
				tableRefKey,
				nextSymbol,
				slr.CollapseEntry{
					RuleNumber:           grammarEntry.RuleNumber,
					Symbol:               grammary.NewSymbol("R"),
					CountOfSymbolsInRule: countOfSymbolsInRule,
				},
			)
			continue
		}

		strategy.recursivelyFindCollapsingEntry(tableRefKey, grammarEntry)
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
		tableRef = strategy.newTableRef()
		strategy.table[tableRefKey][symbol] = tableRef
		strategy.tableRefs = append(strategy.tableRefs, slr.TableEntry{})
		strategy.tableRefsStack = append(strategy.tableRefsStack, tableRef)
	}

	tableEntry := strategy.tableRefs[tableRef]

	tableEntry.GrammarEntries = append(tableEntry.GrammarEntries, grammarEntry)

	if sameTableRef := strategy.fetchSameTableEntry(tableEntry); sameTableRef != nil {
		strategy.table[tableRefKey][symbol] = slr.TableRef(*sameTableRef)

		// mark tableEntry NonValid
		tableEntry.NonValid = true
		strategy.nonValidTableRefs = append(strategy.nonValidTableRefs, tableRef)
		return
	}

	strategy.tableRefs[tableRef] = tableEntry
}

func (strategy *generateStrategy) safeWriteToTableEntryNewCollapseEntry(
	tableRefKey slr.TableRef,
	symbol grammary.Symbol,
	grammarEntry slr.CollapseEntry,
) {
	if _, ok := strategy.table[tableRefKey]; !ok {
		strategy.table[tableRefKey] = map[grammary.Symbol]slr.TableRef{}
	}

	tableRef, ok := strategy.table[tableRefKey][symbol]
	if !ok {
		// if table ref not exists we create new
		tableRef = strategy.newTableRef()
		strategy.table[tableRefKey][symbol] = tableRef
		strategy.tableRefs = append(strategy.tableRefs, slr.TableEntry{})
	}

	tableEntry := strategy.tableRefs[tableRef]

	tableEntry.CollapseEntry = &grammarEntry

	strategy.tableRefs[tableRef] = tableEntry
}

func (strategy *generateStrategy) newTableRef() slr.TableRef {
	if strategy.tableRefsCounter == nil {
		newCounter := slr.TableRef(0)
		strategy.tableRefsCounter = &newCounter
		return *strategy.tableRefsCounter
	}

	*strategy.tableRefsCounter++
	return *strategy.tableRefsCounter
}

func (strategy *generateStrategy) printState() {
	fmt.Println("STATE")

	fmt.Println("STACK", strategy.tableRefsStack)
	fmt.Println("TABLE REFS", strategy.tableRefs)
	fmt.Println("NON VALID TABLE REFS", strategy.nonValidTableRefs)
	fmt.Println("TABLE", strategy.table)

	fmt.Println("END STATE")
}

func (strategy *generateStrategy) findGrammarEntriesForSymbol(symbol grammary.Symbol) []slr.GrammarEntry {
	var result []slr.GrammarEntry
	for ruleNumber, rule := range strategy.grammar.Rules() {
		for numberInRule, symbolInRule := range rule.RuleSymbols() {
			if symbolInRule == symbol {
				result = append(result, slr.GrammarEntry{
					Symbol:       symbolInRule,
					RuleNumber:   uint(ruleNumber),
					NumberInRule: uint(numberInRule),
				})
			}
		}
	}
	return result
}

func (strategy *generateStrategy) fetchSameTableEntry(entry slr.TableEntry) *int {
	for tableRef, tableEntry := range strategy.tableRefs {
		if tableEntry.Equal(entry) {
			return &tableRef
		}
	}
	return nil
}
