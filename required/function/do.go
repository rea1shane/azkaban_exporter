package function

import (
	"azkaban_exporter/pkg/run"
	"azkaban_exporter/required/structs"
)

func Run(e structs.Exporter, errCh chan error) {
	run.Run(e, errCh)
}
