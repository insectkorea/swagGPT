package scanner

import (
	"path/filepath"
	"testing"
)

func TestScanDir(t *testing.T) {
	dir := "testdata"
	files, err := ScanDir(dir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("Expected files, got none")
	}
	if filepath.Base(files[0]) != "example.go" {
		t.Fatalf("Expected example.go, got %s", filepath.Base(files[0]))
	}
}

func TestParseFile(t *testing.T) {
	handlers, err := ParseFile("testdata/example.go")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(handlers) != 2 {
		t.Fatalf("Expected 2 handler functions, got %d", len(handlers))
	}
	expectedHandlers := []string{"Helloworld", "EchoHandler"}
	for i, handler := range handlers {
		if handler.Name.Name != expectedHandlers[i] {
			t.Fatalf("Expected handler %s, got %s", expectedHandlers[i], handler.Name.Name)
		}
	}
}
