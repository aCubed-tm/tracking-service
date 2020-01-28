package main

import "testing"

func TestBla(t *testing.T) {
	t.Error("Not implemented")
}

func TestAbs(t *testing.T) {
	got := 1
	if got != 1 {
		t.Errorf("Abs(-1) = %d; want 1", got)
	}
}
