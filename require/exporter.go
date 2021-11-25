package require

// Exporter basic info
type Exporter struct {
	AppName     string // application name
	TargetName  string // target name for monitoring
	DefaultPort int    // default web listen port
}
