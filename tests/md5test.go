package tests

import (
	"fmt"
	// "main"

	"password-cracker-go/utils"
	"testing"
)

// Function to generate MD5 hash
func GenerateMD5Hash(input string) string {
	return utils.Md5(input)
}

// Test cases for MD5 hash implementation
func TestGenerateMD5Hash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"a", "0cc175b9c0f1b6a831c399e269772661"},
		{"abc", "900150983cd24fb0d6963f7d28e17f72"},
		{"message digest", "f96b697d7cb7938d525a2f31aaf161d0"},
		{"abcdefghijklmnopqrstuvwxyz", "c3fcd3d76192e4007dfb496cca67e13b"},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", "d174ab98d277d9f5a5611c2c9f419d9f"},
		{"12345678901234567890123456789012345678901234567890123456789012345678901234567890", "57edf4a22be3c955ac49da2e2107b67a"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input: %s", tt.input), func(t *testing.T) {
			result := GenerateMD5Hash(tt.input)
			if result != tt.expected {
				t.Errorf("GenerateMD5Hash(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func main() {
	// Run tests
	TestGenerateMD5Hash(&testing.T{})
}
