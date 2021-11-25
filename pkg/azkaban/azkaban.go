package azkaban

type Azkaban struct{}

func (Azkaban) GetTargetName() string {
	return "Azkaban"
}

func (Azkaban) GetAppName() string {
	return "azkaban_exporter"
}

func (Azkaban) GetDefaultListenPort() int {
	return 9900
}
