package common

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"compiler/pkg/common/grammary"
)

var (
	ErrFailedToResolveTableRef = errors.New("failed to resolve table ref")
)

type (
	Table struct {
		AxiomRef          TableRef
		TableMap          TableMap
		TableRefs         TableRefs
		NonValidTableRefs []TableRef
	}

	TableMap map[TableRef]map[grammary.Symbol]TableRef

	TableRef uint

	TableRefs []TableEntry

	// TableEntry has only one of entry not nil else it`s non slr grammar error
	TableEntry struct {
		// grammar entries may not be
		GrammarEntries []GrammarEntry
		// collapse entry may not be
		CollapseEntry *CollapseEntry
		// non valid when we detect same table entry and remove this one
		NonValid bool
	}

	GrammarEntry struct {
		Symbol       grammary.Symbol
		RuleNumber   uint
		NumberInRule uint
	}

	CollapseEntry struct {
		RuleNumber           uint
		Symbol               grammary.Symbol
		CountOfSymbolsInRule uint
	}
)

func (table *Table) ResolveTableRef(ref TableRef) (TableEntry, error) {
	if int(ref) >= len(table.TableRefs) {
		return TableEntry{}, errors.Wrap(ErrFailedToResolveTableRef, fmt.Sprintf("for ref = %d", ref))
	}

	return table.TableRefs[ref], nil
}

func (tableEntry TableEntry) String() string {
	result := make([]string, 0, len(tableEntry.GrammarEntries))
	for _, entry := range tableEntry.GrammarEntries {
		result = append(result, entry.String())
	}

	if tableEntry.CollapseEntry != nil {
		result = append(result, tableEntry.CollapseEntry.String())
	}

	return strings.Join(result, ",")
}

func (tableEntry TableEntry) Equal(other TableEntry) bool {
	if len(tableEntry.GrammarEntries) != len(other.GrammarEntries) {
		return false
	}

	inSlice := func(needle GrammarEntry, haystack []GrammarEntry) bool {
		for _, entry := range haystack {
			if entry == needle {
				return true
			}
		}
		return false
	}

	for _, entry := range tableEntry.GrammarEntries {
		if !inSlice(entry, other.GrammarEntries) {
			return false
		}
	}

	if tableEntry.CollapseEntry != other.CollapseEntry {
		return false
	}

	return true
}

func (grammarEntry GrammarEntry) String() string {
	return fmt.Sprintf(
		"%s:%d:%d",
		grammarEntry.Symbol.String(),
		grammarEntry.RuleNumber,
		grammarEntry.NumberInRule,
	)
}

func (collapseEntry CollapseEntry) String() string {
	return fmt.Sprintf(
		"R:%d",
		collapseEntry.RuleNumber,
	)
}
