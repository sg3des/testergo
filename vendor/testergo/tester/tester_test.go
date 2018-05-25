package tester

import (
	"testing"
	"time"
)

func TestTester(t *testing.T) {
	tester, err := NewTester("./testdata/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err = tester.Start(); err != nil {
		t.Error(err)
	}

	time.Sleep(10 * time.Second)
}
