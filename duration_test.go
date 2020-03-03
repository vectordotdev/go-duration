package duration

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func mustParse(t *testing.T, s string) Duration {
	d, err := ParseRFC3339(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

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

func TestDuration_AddToTime(t *testing.T) {
	tests := []struct {
		t string
		d string

		expected string
	}{
		{
			d: "P3Y6M4DT12H30M5S",
			t: "2006-01-02T15:04:05Z",

			expected: "2009-07-07T03:34:10Z",
		},
		{
			d: "P5W",
			t: "2006-01-02T15:04:05Z",

			expected: "2006-02-06T15:04:05Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			dur := mustParse(t, tt.d)
			tim, err := time.Parse(time.RFC3339, tt.t)
			if err != nil {
				t.Fatal(err)
			}

			expected, err := time.Parse(time.RFC3339, tt.expected)
			if err != nil {
				t.Fatal(err)
			}

			got := dur.AddToTime(tim)

			if !got.Equal(expected) {
				t.Errorf("expected %s.AddToTime(%s) returned %s, expected %s", dur, tim, got, expected)
			}
		})
	}
}

func TestDuration_SubtractFromTime(t *testing.T) {
	tests := []struct {
		t string
		d string

		expected string
	}{
		{
			d: "P3Y6M4DT12H30M5S",
			t: "2006-01-02T15:04:05Z",

			expected: "2002-06-28T02:34:00Z",
		},
		{
			d: "P5W",
			t: "2006-01-02T15:04:05Z",

			expected: "2005-11-28T15:04:05Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			dur := mustParse(t, tt.d)
			tim, err := time.Parse(time.RFC3339, tt.t)
			if err != nil {
				t.Fatal(err)
			}

			expected, err := time.Parse(time.RFC3339, tt.expected)
			if err != nil {
				t.Fatal(err)
			}

			got := dur.SubtractFromTime(tim)

			if !got.Equal(expected) {
				t.Errorf("expected %s.SubtractFromTime(%s) returned %s, expected %s", dur, tim, got, expected)
			}
		})
	}
}
