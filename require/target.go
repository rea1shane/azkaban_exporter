package require

type Target interface {
	GetTargetName() string
	GetAppName() string
	GetDefaultListenPort() int
}
