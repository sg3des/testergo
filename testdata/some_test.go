package name

import "testing"

func TestLog(t *testing.T) {
	t.Log("log message")
}

func TestFail(t *testing.T) {
	t.Error("some error message")
}

func TestPanic(t *testing.T) {
	slice := []string{"string"}
	t.Log(slice)
	t.Log(slice[1])
}
