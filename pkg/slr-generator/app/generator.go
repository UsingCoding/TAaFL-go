package app

import "compiler/pkg/common/grammary"

const newAxiomSymbol = "$F$"

type (
	Table struct {
		Impl         map[grammary.Symbol]map[TableEntryRef]TableEntryRef
		TableEntries []TableEntry
	}

	TableEntryRef uint

	TableEntry struct {
		// Only one field should be not null
		GrammarEntries []GrammarEntry // not nil if table entry is non terminal or terminal with grammar entry
		*CollapseRule                 // not nil if table entry is collapse rule
	}

	GrammarEntry struct {
		referee      grammary.Symbol
		altNumber    uint
		numberInRoll uint
	}

	CollapseRule struct {
		nonTerminal grammary.Symbol
		altNumber   uint
	}
)

func GenerateTable(grammar grammary.Grammar) (Table, error) {
	newGrammar := populateGrammar(grammar)

	table := Table{}

	// Set table header with terminals and non terminals and end of sequence symbol
	for symbol := range newGrammar.Impl {
		table.Impl[symbol] = map[TableEntryRef]TableEntryRef{}
	}
	table.Impl[grammary.NewSymbol(grammary.EndOfSequence)] = map[TableEntryRef]TableEntryRef{}

	err := proceedGrammar(newGrammar, table)

	return table, err
}

func proceedGrammar(grammar grammary.Grammar, table Table) error {

}

func populateGrammar(grammar grammary.Grammar) grammary.Grammar {
	populatedGrammar := grammar.Copy()

	newAxiom := grammary.NewSymbol(newAxiomSymbol)

	populatedGrammar.Impl[newAxiom] = [][]grammary.Symbol{{
		*grammar.Axiom,
	}}

	populatedGrammar.Axiom = &newAxiom

	return populatedGrammar
}
