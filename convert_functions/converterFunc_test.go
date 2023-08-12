package convert_functions

import (
	"testing"
)

func TestConvert_url_to_string(t *testing.T) {
	str, err := Convert_url_to_string("http://example.com")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(str) != ShortURLLength {
		t.Errorf("Unexpected length of converted string. Expected: %d, Got: %d", ShortURLLength, len(str))
	}

	for _, char := range str {
		if !Contains([]byte(AllowedChars), byte(char)) {
			t.Errorf("Unexpected character in converted string: %c", char)
		}
	}
}
