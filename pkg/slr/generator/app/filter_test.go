package app

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"compiler/pkg/common/grammary"
	"compiler/pkg/slr/common/inlinedgrammary"
)

func TestEmptySymbolFilter_Filter(t *testing.T) {
	filter := NewEmptySymbolFilter()

	filteredGrammar, err := filter.Filter(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<S>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("<A>"),
			grammary.NewSymbol("c"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<B>"), []grammary.Symbol{
			grammary.NewSymbol("<A>"),
			grammary.NewSymbol("b"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<B>"), []grammary.Symbol{
			grammary.NewSymbol("e"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<A>"), []grammary.Symbol{
			grammary.NewSymbol("a"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<A>"), []grammary.Symbol{
			grammary.NewSymbol("e"),
		}),
	))
	assert.NoError(t, err)

	sortGrammar(filteredGrammar)
	assert.Equal(t,
		"{{<F>} [{{<S>} [{<A>} {c}]} {{<S>} [{c}]} {{<F>} [{<S>}]} {{<B>} [{<A>} {b}]} {{<B>} [{b}]} {{<A>} [{a}]}]}",
		fmt.Sprint(filteredGrammar),
	)

}

func TestEmptySymbolFilter_FilterRecursive(t *testing.T) {
	filter := NewEmptySymbolFilter()

	filteredGrammar, err := filter.Filter(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<S>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("<A>"),
			grammary.NewSymbol("d"),
			grammary.NewSymbol("<B>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("<A>"),
			grammary.NewSymbol("c"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<B>"), []grammary.Symbol{
			grammary.NewSymbol("<A>"),
			grammary.NewSymbol("b"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<B>"), []grammary.Symbol{
			grammary.NewSymbol("e"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<A>"), []grammary.Symbol{
			grammary.NewSymbol("a"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<A>"), []grammary.Symbol{
			grammary.NewSymbol("e"),
		}),
	))
	assert.NoError(t, err)

	// sort grammar to have predictable sequence of rules
	sortGrammar(filteredGrammar)

	assert.Equal(
		t,
		"{{<F>} [{{<S>} [{<A>} {d} {<B>}]} {{<S>} [{d} {<B>}]} {{<S>} [{<A>} {d}]} {{<S>} [{<A>} {c}]} {{<S>} [{d}]} {{<S>} [{c}]} {{<F>} [{<S>}]} {{<B>} [{<A>} {b}]} {{<B>} [{b}]} {{<A>} [{a}]}]}",
		fmt.Sprint(filteredGrammar),
	)
}

func sortGrammar(grammar inlinedgrammary.Grammar) {
	sort.SliceStable(grammar.Rules(), func(i, j int) bool {
		aRule := grammar.Rules()[i]
		bRule := grammar.Rules()[j]
		a := aRule.LeftSideSymbol()
		b := bRule.LeftSideSymbol()
		if a.String() == b.String() {
			if len(aRule.RuleSymbols()) != len(bRule.RuleSymbols()) {
				if len(aRule.RuleSymbols()) < len(bRule.RuleSymbols()) {
					return false
				}
				return true
			}

			for ai, aSymbol := range aRule.RuleSymbols() {
				bSymbol := bRule.RuleSymbols()[ai]
				if aSymbol.String() < bSymbol.String() {
					return false
				}
				return true
			}
		}

		if a.String() < b.String() {
			return false
		}
		return true
	})
}
