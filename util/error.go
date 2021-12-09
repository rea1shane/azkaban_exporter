package util

var ch = make(chan error)

func GetErrorChannel() chan error {
	return ch
}
