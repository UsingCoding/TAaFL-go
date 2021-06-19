package app

import (
	"compiler/pkg/slr/common"
)

type TableSerializer interface {
	Export(table common.Table) (string, error)
}
