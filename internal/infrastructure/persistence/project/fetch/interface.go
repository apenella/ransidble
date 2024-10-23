package fetch

// SourceCodeFetcher represents the interface for fetching source code components
type SourceCodeFetcher interface {
	Fetch(source string, workingDir string) error
}
