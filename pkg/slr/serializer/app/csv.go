package app

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/slr/common"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

func NewCSVSerializer() TableSerializer {
	return &csvSerializer{}
}

type csvSerializer struct {
}

func newRecord(length int) []string {
	return make([]string, length)
}

func (serializer *csvSerializer) Serialize(table common.Table) (string, error) {
	headersMap := serializer.buildHeaderSymbolsMap(table.TableMap)

	records := make([][]string, 0, len(table.TableRefs)+1)
	records = append(records, serializer.writeHeader(headersMap))

	for tableRefKey, row := range table.TableMap {
		if serializer.isNonValidTableRef(tableRefKey, table.NonValidTableRefs) {
			continue
		}
		record := newRecord(len(headersMap) + 1)
		tableEntryKey, err := table.ResolveTableRef(tableRefKey)
		if err != nil {
			return "", err
		}

		record[0] = tableEntryKey.String()

		for symbol, tableRef := range row {
			position, ok := headersMap[symbol]
			if !ok {
				return "", errors.New(fmt.Sprintf("not found position for header symbol %s", symbol.String()))
			}

			tableEntry, err2 := table.ResolveTableRef(tableRef)
			if err2 != nil {
				return "", err2
			}
			record[position] = tableEntry.String()
		}

		records = append(records, record)
	}

	builder := &strings.Builder{}

	writer := csv.NewWriter(builder)

	err := writer.WriteAll(records)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func (serializer csvSerializer) buildHeaderSymbolsMap(tableMap common.TableMap) map[grammary.Symbol]int {
	headerSymbolsMap := map[grammary.Symbol]int{}

	fetchNextPosition := func() int {
		return len(headerSymbolsMap) + 1
	}

	for _, row := range tableMap {
		for symbol := range row {
			_, ok := headerSymbolsMap[symbol]
			if ok {
				continue
			}

			headerSymbolsMap[symbol] = fetchNextPosition()
		}
	}

	return headerSymbolsMap
}

func (serializer *csvSerializer) writeHeader(headerMap map[grammary.Symbol]int) (record []string) {
	record = newRecord(len(headerMap) + 1)
	record[0] = "TABLE"

	for symbol, i := range headerMap {
		record[i] = symbol.String()
	}

	return record
}

func (serializer *csvSerializer) isNonValidTableRef(ref common.TableRef, nonValidRefs []common.TableRef) bool {
	for _, nonValidRef := range nonValidRefs {
		if nonValidRef == ref {
			return true
		}
	}
	return false
}
