package utils_test

import (
	"reflect"
	"testing"

	"backend-sistem06.com/utils"
)

// Helper function to create a mock FieldLevel for testing
type mockFieldLevel struct {
	value interface{}
}

func (m *mockFieldLevel) Top() reflect.Value {
	return reflect.Value{}
}

func (m *mockFieldLevel) Parent() reflect.Value {
	return reflect.Value{}
}

func (m *mockFieldLevel) Field() reflect.Value {
	return reflect.ValueOf(m.value)
}

func (m *mockFieldLevel) FieldName() string {
	return "testField"
}

func (m *mockFieldLevel) StructFieldName() string {
	return "TestField"
}

func (m *mockFieldLevel) Param() string {
	return ""
}

func (m *mockFieldLevel) GetTag() string {
	return ""
}

func (m *mockFieldLevel) ExtractType(field reflect.Value) (value reflect.Value, kind reflect.Kind, nullable bool) {
	return reflect.Value{}, reflect.String, false
}

func (m *mockFieldLevel) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.String, false
}

func (m *mockFieldLevel) GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.String, false
}

func (m *mockFieldLevel) GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.String, false, false
}

func (m *mockFieldLevel) GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.String, false, false
}

func TestCustomRtRwCodeValidation(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		// Valid
		{
			name:     "Valid RT/RW - 001",
			value:    "001",
			expected: true,
		},
		{
			name:     "Valid RT/RW - 006",
			value:    "006",
			expected: true,
		},
		{
			name:     "Valid RT/RW - 015",
			value:    "015",
			expected: true,
		},
		{
			name:     "Valid RT/RW - 100",
			value:    "100",
			expected: true,
		},
		{
			name:     "Valid RT/RW - 999",
			value:    "999",
			expected: true,
		},
		{
			name:     "Valid RT/RW - 020",
			value:    "020",
			expected: true,
		},

		// Invalid (mengandung huruf)
		{
			name:     "Invalid RT/RW with leading text - abc123",
			value:    "abc123",
			expected: false,
		},
		{
			name:     "Invalid RT/RW with text - RT001",
			value:    "RT001",
			expected: false,
		},

		// Invalid (kosong / tidak 3 digit)
		{
			name:     "Invalid - empty",
			value:    "",
			expected: false,
		},
		{
			name:     "Invalid - only 2 digits",
			value:    "01",
			expected: false,
		},
		{
			name:     "Invalid - 4 digits",
			value:    "1000",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockFieldLevel{value: tt.value}
			res := utils.CustomRtRwCodeValidation(mock)

			if res != tt.expected {
				t.Errorf("Validation failed for '%s': expected %v, got %v",
					tt.value, tt.expected, res)
			}
		})
	}

}

func TestCustomPostalCodeValidation(t *testing.T) {

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		// Valid Postal Code (Indonesia)
		{"Valid - Jakarta 10110", "10110", true},
		{"Valid - Bandung 40115", "40115", true},
		{"Valid - Surabaya 60213", "60213", true},
		{"Valid - Depok 16416", "16416", true},
		{"Valid - Makassar 90111", "90111", true},
		{"Valid - Random valid", "99999", true},

		// Invalid — leading zero
		{"Invalid - leading zero", "01234", false},

		// Invalid — not 5 digits
		{"Invalid - only 4 digits", "1234", false},
		{"Invalid - 6 digits", "123456", false},

		// Invalid — contains letters
		{"Invalid - contains letter", "12A45", false},
		{"Invalid - lowercase letter", "12a45", false},

		// Invalid — contains symbol
		{"Invalid - contains dash", "12-45", false},
		{"Invalid - contains space", "12 45", false},

		// Invalid — empty string
		{"Invalid - empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockFieldLevel{value: tt.value}

			result := utils.CustomPostalCodeValidation(mock)

			if result != tt.expected {
				t.Errorf("Postal code validation failed for '%s': expected %v, got %v",
					tt.value, tt.expected, result)
			}
		})
	}
}
