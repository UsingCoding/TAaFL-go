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

// Added for more detailed error messages
func assertTableRefsEqual(t *testing.T, expected, actual slr.TableRefs) {
	for i, expectedTableEntry := range expected {
		actualTableEntry := actual[i]
		for j, expectedGrammarEntry := range expectedTableEntry.GrammarEntries {
			actualGrammarEntry := actualTableEntry.GrammarEntries[j]
			assert.Equal(t, expectedGrammarEntry, actualGrammarEntry)
		}
		assert.Equal(t, expectedTableEntry.CollapseEntry, actualTableEntry.CollapseEntry)
	}
}
