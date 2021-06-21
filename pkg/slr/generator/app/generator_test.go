package app

import (
	"fmt"
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
				RuleNumber: 0,
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

func TestGenerator_GenerateTable_RecursiveGrammar(t *testing.T) {
	generator := NewGenerator()
	validator := NewValidator()

	table, err := generator.GenerateTable(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<S>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("("),
			grammary.NewSymbol("<S>"),
			grammary.NewSymbol(")"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("("),
			grammary.NewSymbol(")"),
		}),
	))
	assert.NoError(t, err)

	assert.NoError(t, validator.Validate(table))

	expectedMap := "map[0:map[(:2 <S>:1] 1:map[_|_:13] 2:map[(:2 ):5 <S>:3] 3:map[):10] 4:map[(:2 <S>:3] 5:map[):7 _|_:6] 10:map[):12 _|_:11]]"
	expectedRefs := "[<F>00 <S>00 (10,(20 <S>11 (10 )21 R:2 R:2   )12 R:1 R:1 R:0]"
	expectedNonValidRefs := "[4 8 9 4]"

	assert.Equal(t, expectedMap, fmt.Sprint(table.TableMap))
	assert.Equal(t, expectedRefs, fmt.Sprint(table.TableRefs))
	assert.Equal(t, expectedNonValidRefs, fmt.Sprint(table.NonValidTableRefs))
}

func TestGenerator_GenerateTable_NonValidGrammar(t *testing.T) {
	generator := NewGenerator()
	validator := NewValidator()

	table, err := generator.GenerateTable(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<E>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<E>"), []grammary.Symbol{
			grammary.NewSymbol("<E>"),
			grammary.NewSymbol("+"),
			grammary.NewSymbol("<E>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<E>"), []grammary.Symbol{
			grammary.NewSymbol("i"),
		}),
	))
	assert.NoError(t, err)

	assert.Error(t, validator.Validate(table), ErrNotSlrGrammar{})

	assert.Equal(t, fmt.Sprint(table.TableMap), "map[0:map[<E>:1 i:2] 1:map[+:6 _|_:5] 2:map[+:4 _|_:3] 6:map[<E>:7 i:2] 7:map[+:10 _|_:9]]")
	assert.Equal(t, fmt.Sprint(table.TableRefs), "[<F>00 <E>00,<E>10 i20 R:1 R:1 R:0 +11 <E>12,<E>10  R:1 +11,R:1]")
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
