package require

// Exporter customize exporter's basic info
type Exporter interface {
	// GetName return application name
	GetName() string

	// GetMonitorTargetName return target name for monitoring
	GetMonitorTargetName() string

	// GetDefaultPort return default web listen port
	GetDefaultPort() int
}
