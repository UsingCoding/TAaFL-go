package executor

const (
	nodeJS = "node"
)

func NewNodeJSExecutor() (NodeJS, error) {
	return New(nodeJS)
}

type NodeJS Executor
