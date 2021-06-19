package common

import (
	"compiler/pkg/common/grammary"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

var (
	ErrFailedToResolveTableRef = errors.New("failed to resolve table ref")
)

type (
	Table struct {
		TableMap  TableMap
		TableRefs TableRefs
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
	}

	GrammarEntry struct {
		Symbol       grammary.Symbol
		RuleNumber   uint
		NumberInRule uint
	}

	CollapseEntry struct {
		RuleNumber uint
		//deprecated
		Symbol grammary.Symbol
		//deprecated
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

func (grammarEntry GrammarEntry) String() string {
	return fmt.Sprintf(
		"%s%d%d",
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
