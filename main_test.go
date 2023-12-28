package main_test

import (
	"api"
	"testing"
)

func TestRun(t *testing.T) {
	if err := main.Run(); err != nil {
		t.Fatal(err)
	}
}
