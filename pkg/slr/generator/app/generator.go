package app

import (
	"fmt"
	"sort"

	"compiler/pkg/common/grammary"
	slr "compiler/pkg/slr/common"
	"compiler/pkg/slr/common/inlinedgrammary"
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

				// If symbol is nonTerminal need to proceed transit closure
				if grammarEntry.Symbol.NonTerminal() {
					symbol := grammarEntry.Symbol

					// This additional transit closure adds records which doesn't affect on runner
					strategy.proceedRecursiveTransitClosure(tableRef, symbol, nil)
					continue
				}

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

			newTableRef := strategy.safeWriteToTableEntryNewGrammarEntry(
				tableRef,
				symbol,
				slr.GrammarEntry{
					Symbol:       symbol,
					RuleNumber:   grammarEntry.RuleNumber,
					NumberInRule: nextNumberInRule,
				},
			)

			strategy.proceedRecursiveTransitClosure(tableRef, symbol, &newTableRef)
		}
	}

	return slr.Table{
		AxiomRef:          slr.TableRef(0),
		TableMap:          strategy.table,
		TableRefs:         strategy.tableRefs,
		NonValidTableRefs: strategy.nonValidTableRefs,
	}, nil
}

func (strategy *generateStrategy) proceedRecursiveTransitClosure(
	tableRefKey slr.TableRef,
	handledNonTerminal grammary.Symbol,
	handledNonTerminalTableRef *slr.TableRef,
) {
	strategy.printState()

	if handledNonTerminalTableRef != nil {
		strategy.proceedRecursiveTransitClosure(*handledNonTerminalTableRef, handledNonTerminal, nil)
	}

	rulesMap := strategy.grammar.FindByLeftSideSymbol(handledNonTerminal)

	keys := make([]int, 0, len(rulesMap))
	for k := range rulesMap {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	// we sort map keys to have predictable order in grammar entry
	for _, key := range keys {
		ruleNumber := uint(key)
		rule := rulesMap[ruleNumber]
		const firstSymbolPos = 0
		symbol := rule.RuleSymbols()[firstSymbolPos]

		if symbol.NonTerminal() {
			newTableRef := strategy.safeWriteToTableEntryNewGrammarEntry(
				tableRefKey,
				symbol,
				slr.GrammarEntry{
					Symbol:       symbol,
					RuleNumber:   ruleNumber,
					NumberInRule: firstSymbolPos,
				},
			)

			strategy.proceedRecursiveTransitClosure(tableRefKey, symbol, &newTableRef)

			continue
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
		if entry == grammarEntry {
			continue
		}

		rule := strategy.grammar.Rules()[entry.RuleNumber]
		countOfSymbolsInRule := uint(len(rule.RuleSymbols()))

		// checking that element last in axiom rule
		if rule.LeftSideSymbol() == strategy.grammar.Axiom() && entry.NumberInRule+1 == countOfSymbolsInRule {
			strategy.safeWriteToTableEntryNewCollapseEntry(
				tableRefKey,
				grammary.NewSymbol(grammary.EndOfSequence),
				slr.CollapseEntry{
					RuleNumber:           grammarEntry.RuleNumber,
					Symbol:               symbolCollapsingTo,
					CountOfSymbolsInRule: countOfSymbolsInCollapsingRule,
				},
			)
			continue
		}

		if entry.NumberInRule+1 >= uint(len(rule.RuleSymbols())) {
			strategy.recursivelyFindCollapsingEntry(tableRefKey, entry)
			continue
		}

		nextSymbol := rule.RuleSymbols()[entry.NumberInRule+1]

		if !nextSymbol.NonTerminal() {
			strategy.safeWriteToTableEntryNewCollapseEntry(
				tableRefKey,
				nextSymbol,
				slr.CollapseEntry{
					RuleNumber:           grammarEntry.RuleNumber,
					Symbol:               symbolCollapsingTo,
					CountOfSymbolsInRule: countOfSymbolsInCollapsingRule,
				},
			)
			continue
		}

		strategy.findCollapsingEntryForNonTerminalViaTransitClosure(
			tableRefKey,
			symbolCollapsingTo,
			countOfSymbolsInCollapsingRule,
			slr.GrammarEntry{
				Symbol:       nextSymbol,
				RuleNumber:   entry.RuleNumber,
				NumberInRule: entry.NumberInRule + 1,
			},
		)
	}
}

// More of func arguments just passed to be recorded to table via safeWriteToTableEntryNewCollapseEntry
func (strategy *generateStrategy) findCollapsingEntryForNonTerminalViaTransitClosure(
	tableRefKey slr.TableRef,
	symbolCollapsingTo grammary.Symbol,
	countOfSymbolsInCollapsingRule uint,
	grammarEntry slr.GrammarEntry,
) {
	rulesMap := strategy.grammar.FindByLeftSideSymbol(grammarEntry.Symbol)

	keys := make([]int, 0, len(rulesMap))
	for k := range rulesMap {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	// we sort map keys to have predictable order in grammar entry
	for _, key := range keys {
		ruleNumber := uint(key)
		rule := rulesMap[ruleNumber]
		const firstSymbolPos = 0
		symbol := rule.RuleSymbols()[firstSymbolPos]

		if symbol.NonTerminal() {
			strategy.findCollapsingEntryForNonTerminalViaTransitClosure(
				tableRefKey,
				symbolCollapsingTo,
				countOfSymbolsInCollapsingRule,
				slr.GrammarEntry{
					Symbol:       symbol,
					RuleNumber:   ruleNumber,
					NumberInRule: firstSymbolPos,
				},
			)
			continue
		}

		strategy.safeWriteToTableEntryNewCollapseEntry(
			tableRefKey,
			symbol,
			slr.CollapseEntry{
				RuleNumber:           grammarEntry.RuleNumber,
				Symbol:               symbolCollapsingTo,
				CountOfSymbolsInRule: countOfSymbolsInCollapsingRule,
			},
		)
	}
}

func (strategy *generateStrategy) findCollapsingEntryForNonTerminalViaTransitClosureDeprecated(
	tableRefKey slr.TableRef,
	symbolCollapsingTo grammary.Symbol,
	countOfSymbolsInCollapsingRule uint,
	grammarEntry slr.GrammarEntry,
) {
	rule := strategy.grammar.Rules()[grammarEntry.RuleNumber]
	nextSymbolNumberInRule := grammarEntry.NumberInRule + 1

	if nextSymbolNumberInRule == uint(len(rule.RuleSymbols())) {
		panic("incorrect usage of generateStrategy.findCollapsingEntryForNonTerminalViaTransitClosure: rule is out of symbols")
	}

	nextSymbol := rule.RuleSymbols()[nextSymbolNumberInRule]

	newGrammarEntry := slr.GrammarEntry{
		Symbol:       nextSymbol,
		RuleNumber:   grammarEntry.RuleNumber,
		NumberInRule: nextSymbolNumberInRule,
	}

	if nextSymbol.NonTerminal() {
		strategy.findCollapsingEntryForNonTerminalViaTransitClosure(
			tableRefKey,
			symbolCollapsingTo,
			countOfSymbolsInCollapsingRule,
			newGrammarEntry,
		)
		return
	}

	strategy.recursivelyFindCollapsingEntry(tableRefKey, grammarEntry)
}

func (strategy *generateStrategy) safeWriteToTableEntryNewGrammarEntry(
	tableRefKey slr.TableRef,
	symbol grammary.Symbol,
	grammarEntry slr.GrammarEntry,
) slr.TableRef {
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

	// If in table entry already exists grammar entry don`t allow adding duplicate
	for _, entry := range tableEntry.GrammarEntries {
		if entry == grammarEntry {
			return tableRef
		}
	}

	tableEntry.GrammarEntries = append(tableEntry.GrammarEntries, grammarEntry)

	if sameTableRef := strategy.fetchSameTableEntry(tableEntry); sameTableRef != nil {
		strategy.table[tableRefKey][symbol] = slr.TableRef(*sameTableRef)

		// mark tableEntry NonValid
		tableEntry.NonValid = true
		strategy.nonValidTableRefs = append(strategy.nonValidTableRefs, tableRef)
		return slr.TableRef(*sameTableRef)
	}

	strategy.tableRefs[tableRef] = tableEntry
	return tableRef
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
