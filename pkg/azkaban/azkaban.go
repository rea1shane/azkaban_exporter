package azkaban

type Azkaban struct{}

func (Azkaban) GetName() string {
	return "azkaban_exporter"
}

func (Azkaban) GetMonitorTargetName() string {
	return "Azkaban"
}

func (Azkaban) GetDefaultPort() int {
	return 9900
}
