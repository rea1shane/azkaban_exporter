package required

type Target interface {
	// GetNamespace return the pre of metrics
	GetNamespace() string
}
