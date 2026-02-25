package verify

import "testing"

func TestIsValidSHA256HexLower(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want bool
	}{
		{
			name: "valid",
			in:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			want: true,
		},
		{
			name: "uppercase_false",
			in:   "0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF",
			want: false,
		},
		{
			name: "wrong_length_false",
			in:   "deadbeef",
			want: false,
		},
		{
			name: "invalid_char_false",
			in:   "g123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			want: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := IsValidSHA256HexLower(tc.in); got != tc.want {
				t.Fatalf("IsValidSHA256HexLower(%q)=%v want=%v", tc.in, got, tc.want)
			}
		})
	}
}
