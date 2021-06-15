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

	TableEntry []GrammarEntry

	GrammarEntry struct {
		Symbol       grammary.Symbol
		RuleNumber   uint
		NumberInRule uint
	}
)

func (table *Table) ResolveTableRef(ref TableRef) (TableEntry, error) {
	if int(ref) >= len(table.TableRefs) {
		return nil, errors.Wrap(ErrFailedToResolveTableRef, fmt.Sprintf("for ref = %d", ref))
	}

	return table.TableRefs[ref], nil
}

func (tableEntry TableEntry) String() string {
	result := make([]string, 0, len(tableEntry))
	for _, entry := range tableEntry {
		result = append(result, fmt.Sprintf(
			"%s%d%d",
			entry.Symbol.String(),
			entry.RuleNumber,
			entry.NumberInRule,
		))
	}
	return strings.Join(result, ",")
}
