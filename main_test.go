package main

import (
	"testing"
)

func TestComputeSHA1(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple json",
			input:    `{"key":"value"}`,
			expected: "228458095a9502070fc113d99504226a6ff90a9a",
		},
		{
			name:     "empty json",
			input:    `{}`,
			expected: "bf21a9e8fbc5a3846fb05b4fa0859e0917b2202f",
		},
		{
			name:     "complex json",
			input:    `{"user":"john","age":30,"active":true}`,
			expected: "99811ccac5178656a48f20f9f24656c88ee2a0de",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := computeSHA1([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("computeSHA1(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestComputeSHA1Deterministic(t *testing.T) {
	// Test that the same input always produces the same hash
	input := `{"test":"data","value":123}`
	hash1 := computeSHA1([]byte(input))
	hash2 := computeSHA1([]byte(input))
	
	if hash1 != hash2 {
		t.Errorf("computeSHA1 is not deterministic: first call = %q, second call = %q", hash1, hash2)
	}
}

func TestComputeSHA1Different(t *testing.T) {
	// Test that different inputs produce different hashes
	input1 := `{"key":"value1"}`
	input2 := `{"key":"value2"}`
	
	hash1 := computeSHA1([]byte(input1))
	hash2 := computeSHA1([]byte(input2))
	
	if hash1 == hash2 {
		t.Errorf("computeSHA1 produced the same hash for different inputs: %q", hash1)
	}
}
