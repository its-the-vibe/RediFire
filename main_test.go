package main

import (
	"testing"
)

func TestComputeSHA256(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple json",
			input:    `{"key":"value"}`,
			expected: "e43abcf3375244839c012f9633f95862d232a95b00d5bc7348b3098b9fed7f32",
		},
		{
			name:     "empty json",
			input:    `{}`,
			expected: "44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
		},
		{
			name:     "complex json",
			input:    `{"user":"john","age":30,"active":true}`,
			expected: "906a0da0aac66ac39c4c8f2174f82542d7b732adceb2b8eb38632e13e30848ae",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := computeSHA256(tt.input)
			if result != tt.expected {
				t.Errorf("computeSHA256(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestComputeSHA256Deterministic(t *testing.T) {
	// Test that the same input always produces the same hash
	input := `{"test":"data","value":123}`
	hash1 := computeSHA256(input)
	hash2 := computeSHA256(input)
	
	if hash1 != hash2 {
		t.Errorf("computeSHA256 is not deterministic: first call = %q, second call = %q", hash1, hash2)
	}
}

func TestComputeSHA256Different(t *testing.T) {
	// Test that different inputs produce different hashes
	input1 := `{"key":"value1"}`
	input2 := `{"key":"value2"}`
	
	hash1 := computeSHA256(input1)
	hash2 := computeSHA256(input2)
	
	if hash1 == hash2 {
		t.Errorf("computeSHA256 produced the same hash for different inputs: %q", hash1)
	}
}
