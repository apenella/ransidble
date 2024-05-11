package executor

// Executor represents an executor to run tasks
type Executor interface {
	Run() error
}
