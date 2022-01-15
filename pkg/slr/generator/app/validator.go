package app

import (
	"fmt"

	"compiler/pkg/slr/common"
)

type ErrNotSlrGrammar struct {
	tableRef   int
	tableEntry common.TableEntry
}

func (e *ErrNotSlrGrammar) Error() string {
	return fmt.Sprintf("not slr grammar: table ref #%d has both grammar entries and collapse rule: %s", e.tableRef, e.tableEntry)
}

type Validator interface {
	// Validate validates table
	Validate(table common.Table) error
}

func NewValidator() Validator {
	return &validator{}
}

type validator struct {
}

func (v *validator) Validate(table common.Table) error {
	for tableRef, tableEntry := range table.TableRefs {
		if tableEntry.GrammarEntries != nil && tableEntry.CollapseEntry != nil {
			return &ErrNotSlrGrammar{
				tableRef:   tableRef,
				tableEntry: tableEntry,
			}
		}
	}

	return nil
}
