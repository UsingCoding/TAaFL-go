package slr

import (
	"compiler/pkg/slr/common/inlinedgrammary"
	"compiler/pkg/slr/export"
	"compiler/pkg/slr/generator/app"
)

type Analyzer interface {
	Analyze(grammar inlinedgrammary.Grammar) error
}

func NewAnalyzer(
	tableGenerator app.Generator,
	tableValidator app.Validator,
	tableExporter export.TableExporter,
) Analyzer {
	return &analyzer{
		tableGenerator: tableGenerator,
		tableValidator: tableValidator,
		tableExporter:  tableExporter,
	}
}

type analyzer struct {
	tableGenerator app.Generator
	tableValidator app.Validator
	tableExporter  export.TableExporter
}

func (slr *analyzer) Analyze(grammar inlinedgrammary.Grammar) error {
	table, err := slr.tableGenerator.GenerateTable(grammar)
	if err != nil {
		return err
	}

	err = slr.tableValidator.Validate(table)
	if err != nil {
		return err
	}

	err = slr.tableExporter.Export(table)
	if err != nil {
		return err
	}

	return nil
}
