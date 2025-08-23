package utils

import (
	"os"
	"testing"
)

func TestReadStringInto(t *testing.T) {
	// Prepare input and redirect os.Stdin
	input := "hello world\n"
	r, w, _ := os.Pipe()
	w.WriteString(input)
	w.Close()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	var result string
	ReadStringInto("Enter: ", &result)

	if result != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", result)
	}
}

func TestReadStringSliceInto(t *testing.T) {
	input := "a, b, c ,, d\n"
	r, w, _ := os.Pipe()
	w.WriteString(input)
	w.Close()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	var result []string
	ReadStringSliceInto("Enter: ", &result)

	expected := []string{"a", "b", "c", "d"}
	if len(result) != len(expected) {
		t.Fatalf("Expected %d elements, got %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("At index %d: expected '%s', got '%s'", i, v, result[i])
		}
	}
}
