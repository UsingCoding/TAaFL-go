package executor

const (
	script = "/app/bin/astbackend"
)

func NewASTBackendExecutor() (ASTBackendExecutor, error) {
	return New(script)
}

type ASTBackendExecutor Executor
