package app

import "encoding/json"

// Builder builds AST representation
type Builder interface {
	Build(stack Stack) (json.RawMessage, error)
}
