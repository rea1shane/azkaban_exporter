package azkaban

type Azkaban struct {
	Namespace string
	Address   []string
}

func (a Azkaban) GetNamespace() string {
	return a.Namespace
}
