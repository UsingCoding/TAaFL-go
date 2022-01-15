package app

import (
	"compiler/pkg/slr/common"
)

type TableSerializer interface {
	Serialize(table common.Table) (string, error)
}
