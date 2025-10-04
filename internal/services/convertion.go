package services

import (
	"strings"
	"unicode"
)

// this function converts any input string like "kamaraj college"
func Convert_to_CamelCase(s string) string {

	// Remove leading and trailing spaces from input
	s = strings.TrimSpace(s)
	parts := strings.Fields(strings.ToLower(s)) // Convert entire string to lowercase, then split by spaces into words

	var b strings.Builder // Create a string builder to efficiently concatenate strings
	for _, part := range parts {
		if len(part) > 0 {
			// Capitalize first letter of the word and append the rest of the word as-is
			b.WriteString(strings.ToUpper(string(part[0])) + part[1:])
		}
	}
	return b.String()
}

// this function converts a CamelCase string like "KamarajCollege"
// back to a human-readable format like "Kamaraj College"
func Convert_from_CamelCase(s string) string {
	var b strings.Builder // Use string builder for efficient concatenation
	for i, r := range s {

		// If the character is uppercase and not the first character
		if i > 0 && unicode.IsUpper(r) {
			b.WriteByte(' ') // Insert a space before it
		}
		b.WriteRune(r) // Append the character itself
	}
	return b.String()
}
