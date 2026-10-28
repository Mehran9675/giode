package main

import (
	"bytes"
	"os"
	"testing"
)

// TestUpToDate fails when reexport.go wasn't regenerated after package
// ui changed: run go generate in the repository root.
func TestUpToDate(t *testing.T) {
	want, err := generate("../../../ui")
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile("../../../reexport.go")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Error("reexport.go is stale: run go generate in the repository root")
	}
}
