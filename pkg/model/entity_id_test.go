package model

import "testing"

func TestIsEntityID(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  bool
	}{
		{"54013ba69c196820e56801f1ef5aad54", true},
		{"cart", false},
		{"54013BA69C196820E56801F1EF5AAD54", false},
		{"550e8400-e29b-41d4-a716-446655440000", false},
		{"54013ba69c196820e56801f1ef5aad5", false},
	} {
		if got := IsEntityID(tc.value); got != tc.want {
			t.Fatalf("IsEntityID(%q)=%v, want %v", tc.value, got, tc.want)
		}
	}
}
