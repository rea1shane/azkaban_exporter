package structs

// Exporter basic info
type Exporter struct {
	MonitorTargetName string // MonitorTargetName is target name for monitoring, will convert to exporter.Exporter's Namespace and ExporterName
	DefaultPort       int    // DefaultPort is default web listen port of exporter
}
