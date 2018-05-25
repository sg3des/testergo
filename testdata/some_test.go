package name

import "testing"

func TestLog(t *testing.T) {
	t.Log("log")
}

func TestFail(t *testing.T) {
	t.Error("asdasd")
}

func TestPanic(t *testing.T) {
	slice := []string{"string"}
	t.Log(slice)
	t.Log(slice[1])
}
