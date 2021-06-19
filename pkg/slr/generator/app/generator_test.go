package app

import (
	"testing"

	slr "compiler/pkg/slr/common"
	"github.com/stretchr/testify/assert"

	"compiler/pkg/common/grammary"
	"compiler/pkg/slr/common/inlinedgrammary"
)

func TestGenerator_GenerateTable_SimpleGrammar(t *testing.T) {
	generator := NewGenerator()
	validator := NewValidator()

	table, err := generator.GenerateTable(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("a"),
		}),
	))
	assert.NoError(t, err)

	assert.NoError(t, validator.Validate(table))

	assert.Equal(t, slr.TableMap(map[slr.TableRef]map[grammary.Symbol]slr.TableRef{
		slr.TableRef(0): {
			grammary.NewSymbol("a"): slr.TableRef(1),
		},
		slr.TableRef(1): {
			grammary.NewSymbol(grammary.EndOfSequence): slr.TableRef(2),
		},
	}),
		table.TableMap,
	)

	assertTableRefsEqual(t, slr.TableRefs{
		{
			GrammarEntries: []slr.GrammarEntry{
				{
					Symbol:       grammary.NewSymbol("<F>"),
					RuleNumber:   0,
					NumberInRule: 0,
				},
			},
		},
		{
			GrammarEntries: []slr.GrammarEntry{
				{
					Symbol:       grammary.NewSymbol("a"),
					RuleNumber:   0,
					NumberInRule: 0,
				},
			},
		},
		{
			CollapseEntry: &slr.CollapseEntry{
				RuleNumber:           0,
				Symbol:               grammary.NewSymbol("R"),
				CountOfSymbolsInRule: 1,
			},
		},
	},
		table.TableRefs,
	)

	assert.Equal(t, 0, len(table.NonValidTableRefs))
}

func TestGenerator_GenerateTable_SecondGrammar(t *testing.T) {
	generator := NewGenerator()
	validator := NewValidator()

	table, err := generator.GenerateTable(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<S>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("a"),
			grammary.NewSymbol("b"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("c"),
		}),
	))
	assert.NoError(t, err)

	assert.NoError(t, validator.Validate(table))

	assert.Equal(t, slr.TableMap(map[slr.TableRef]map[grammary.Symbol]slr.TableRef{
		slr.TableRef(0): {
			grammary.NewSymbol("<S>"): slr.TableRef(1),
			grammary.NewSymbol("a"):   slr.TableRef(2),
			grammary.NewSymbol("c"):   slr.TableRef(3),
		},
		slr.TableRef(1): {
			grammary.NewSymbol(grammary.EndOfSequence): slr.TableRef(7),
		},
		slr.TableRef(2): {
			grammary.NewSymbol("b"): slr.TableRef(5),
		},
		slr.TableRef(3): {
			grammary.NewSymbol(grammary.EndOfSequence): slr.TableRef(4),
		},
		slr.TableRef(5): {
			grammary.NewSymbol(grammary.EndOfSequence): slr.TableRef(6),
		},
	}),
		table.TableMap,
	)

	assertTableRefsEqual(t, slr.TableRefs{
		{
			GrammarEntries: []slr.GrammarEntry{
				{
					Symbol:       grammary.NewSymbol("<F>"),
					RuleNumber:   0,
					NumberInRule: 0,
				},
			},
		},
		{
			GrammarEntries: []slr.GrammarEntry{
				{
					Symbol:       grammary.NewSymbol("<S>"),
					RuleNumber:   0,
					NumberInRule: 0,
				},
			},
		},
		{
			GrammarEntries: []slr.GrammarEntry{
				{
					Symbol:       grammary.NewSymbol("a"),
					RuleNumber:   1,
					NumberInRule: 0,
				},
			},
		},
		{
			GrammarEntries: []slr.GrammarEntry{
				{
					Symbol:       grammary.NewSymbol("c"),
					RuleNumber:   2,
					NumberInRule: 0,
				},
			},
		},
		{
			CollapseEntry: &slr.CollapseEntry{
				RuleNumber: 2,
			},
		},
		{
			GrammarEntries: []slr.GrammarEntry{
				{
					Symbol:       grammary.NewSymbol("b"),
					RuleNumber:   1,
					NumberInRule: 1,
				},
			},
		},
		{
			CollapseEntry: &slr.CollapseEntry{
				RuleNumber: 1,
			},
		},
		{
			CollapseEntry: &slr.CollapseEntry{
				RuleNumber: 0,
			},
		},
	},
		table.TableRefs,
	)
}

// Added for more detailed error messages
func assertTableRefsEqual(t *testing.T, expected, actual slr.TableRefs) {
	for i, expectedTableEntry := range expected {
		actualTableEntry := actual[i]
		for j, expectedGrammarEntry := range expectedTableEntry.GrammarEntries {
			actualGrammarEntry := actualTableEntry.GrammarEntries[j]
			assert.Equal(t, expectedGrammarEntry, actualGrammarEntry)
		}
		if expectedTableEntry.CollapseEntry == nil {
			assert.Nil(t, actualTableEntry.CollapseEntry)
			continue
		}

		assert.Equal(t, expectedTableEntry.CollapseEntry.RuleNumber, actualTableEntry.CollapseEntry.RuleNumber)
	}
}
