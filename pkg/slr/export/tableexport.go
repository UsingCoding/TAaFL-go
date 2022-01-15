package export

import (
	"compiler/pkg/slr/common"
	"compiler/pkg/slr/serializer/app"
	"os"
)

type TableExporter interface {
	Export(table common.Table) error
}

func NewFileTableExporter(serializer app.TableSerializer, path string) TableExporter {
	return &fileTableExporter{serializer: serializer, path: path}
}

type fileTableExporter struct {
	serializer app.TableSerializer
	path       string
}

func (exporter *fileTableExporter) Export(table common.Table) error {
	serializedTable, err := exporter.serializer.Serialize(table)
	if err != nil {
		return err
	}

	file, err := os.Create(exporter.path)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()

	_, err = file.WriteString(serializedTable)
	if err != nil {
		return err
	}

	return nil
}

func NewOptionalExporter(exporter TableExporter, export bool) TableExporter {
	return &optionalExporter{
		exporter: exporter,
		export:   export,
	}
}

type optionalExporter struct {
	exporter TableExporter
	export   bool
}

func (decorator *optionalExporter) Export(table common.Table) error {
	if !decorator.export {
		return nil
	}

	return decorator.exporter.Export(table)
}
