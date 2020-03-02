package duration

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseRFC3339(t *testing.T) {
	tests := []struct {
		s string

		expected    Duration
		expectedErr error
	}{
		{
			s: "P3Y6M4DT12H30M5S",

			expected: Duration{
				Years:   3,
				Months:  6,
				Days:    4,
				Hours:   12,
				Minutes: 30,
				Seconds: 5,
			},
		},
		{
			s: "P5W",

			expected: Duration{
				Weeks: 5,
			},
		},
		{
			s: "P3Y5W",

			expectedErr: errors.New("must be RFC3999 formatted duration"),
		},
		{
			s: "123",

			expectedErr: errors.New("must be RFC3999 formatted duration"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got, gotErr := ParseRFC3339(tt.s)

			switch {
			case gotErr == nil && tt.expectedErr != nil:
				t.Errorf("ParseRFC3339(%q) returned no error, expected error %q", tt.s, tt.expectedErr)
				return
			case gotErr != nil && tt.expectedErr == nil:
				t.Errorf("ParseRFC3339(%q) returned error %q, expected no error", tt.s, gotErr)
				return
			case gotErr != nil && tt.expectedErr != nil:
				if tt.expectedErr.Error() != gotErr.Error() {
					t.Errorf("ParseRFC3339(%q) returned error %q, expected error %q", tt.s, gotErr, tt.expectedErr)
				}
			}

			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("expected ParseRFC3339(%q) mismatch (-want +got):\n%s", tt.s, diff)
			}
		})
	}
}

func TestDuration_FormatRFC3339(t *testing.T) {
	tests := []struct {
		d Duration

		expected string
	}{
		{
			d: Duration{
				Years:   3,
				Months:  6,
				Days:    4,
				Hours:   12,
				Minutes: 30,
				Seconds: 5,
			},

			expected: "P3Y6M4DT12H30M5S",
		},
		{
			d: Duration{
				Weeks: 5,
			},

			expected: "P5W",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := tt.d.FormatRFC3339()

			if got != tt.expected {
				t.Errorf("expected %+v.FormatRFC3339() returned %s, expected %s", tt.d, got, tt.expected)
			}
		})
	}
}
