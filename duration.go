/*
* Package duration allows for parsing and use of RFC3339 duration values
 */
package duration

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// ABNF pattern in https://www.ietf.org/rfc/rfc3339.txt
var (
	rfc3339DurationPattern = regexp.MustCompile(`\AP(` +
		// a number of weeks
		`((?P<weeks>\d+)W)` +
		`|(` +
		// date duration
		`((?P<years>\d+)Y)?((?P<months>\d+)M)?((?P<days>\d+)D)?` +
		`(T((?P<hours>\d+)H)?((?P<minutes>\d+)M)?((?P<seconds>\d+?)S)?)?` +
		`)` +
		`)\z`,
	)

	ErrInvalidFormat = fmt.Errorf("must be RFC3999 formatted duration")
)

type ParseError struct {
	error
}

// Duration represents a time duration
//
// We do not use time.Duration as there would be ambiguilty when units greater than a day are used. For example, if we
// encoded a year as 365 days, adding the duration to a date would not yield the correct result when a leap year is in
// the interval.
type Duration struct {
	Years  int
	Months int

	// TODO(jesse):
	// Weeks should not be combined with other units according to RFC 3339 we could seperate this into two structs with an
	// interface (like WeekDuration)
	Weeks int

	Days    int
	Hours   int
	Minutes int
	Seconds int
}

// ParseRFC3339 parses a duration encdoded as described in RFC 3339
func ParseRFC3339(s string) (Duration, error) {
	fmt.Println(rfc3339DurationPattern.MatchString(s))
	matches := rfc3339DurationPattern.FindStringSubmatch(s)
	if matches == nil {
		return Duration{}, ParseError{fmt.Errorf("must be RFC3999 formatted duration")}
	}

	d := Duration{}
	for i, name := range rfc3339DurationPattern.SubexpNames() {
		value := matches[i]
		if i == 0 || name == "" || len(value) == 0 {
			continue
		}

		i64, err := strconv.ParseInt(string(value), 10, 32)
		if err != nil {
			return Duration{}, ParseError{fmt.Errorf("must be RFC3999 formatted duration, found non-integer: %s", string(value))}
		}

		i := int(i64)

		switch name {
		case "years":
			d.Years = i
		case "months":
			d.Months = i
		case "weeks":
			d.Weeks = i
		case "days":
			d.Days = i
		case "hours":
			d.Hours = i
		case "minutes":
			d.Minutes = i
		case "seconds":
			d.Seconds = i
		}
	}

	return d, nil
}

// FormatRFC3339 returns the duration formatted as an RFC 3339 string
func (d Duration) FormatRFC3339() string {
	buf := &bytes.Buffer{}
	fmt.Fprint(buf, "P")
	if d.Years != 0 {
		fmt.Fprintf(buf, "%dY", d.Years)
	}
	if d.Months != 0 {
		fmt.Fprintf(buf, "%dM", d.Months)
	}
	if d.Weeks != 0 {
		fmt.Fprintf(buf, "%dW", d.Weeks)
	}
	if d.Days != 0 {
		fmt.Fprintf(buf, "%dD", d.Days)
	}
	if d.Hours != 0 || d.Minutes != 0 || d.Seconds != 0 {
		fmt.Fprint(buf, "T")
	}
	if d.Hours != 0 {
		fmt.Fprintf(buf, "%dH", d.Hours)
	}
	if d.Minutes != 0 {
		fmt.Fprintf(buf, "%dM", d.Minutes)
	}
	if d.Seconds != 0 {
		fmt.Fprintf(buf, "%dS", d.Seconds)
	}
	return buf.String()
}

// String retrurns the duration as an RFC 3339 duration
func (d Duration) String() string {
	return d.FormatRFC3339()
}

// MarshalText implements the encoding.TextMarshaler interface. The duration is formatted in RFC 3339 format.
func (d Duration) MarshalText() (text []byte, err error) {
	return []byte(d.FormatRFC3339()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface. The duration is expected to be in RFC 3339 format.
func (d *Duration) UnmarshalText(text []byte) error {
	duration, err := ParseRFC3339(string(text))
	if err != nil {
		return err
	}

	*d = duration

	return nil
}

// AddTo adds the duration to the provided time.Time, returning a new time.Time.
func (d Duration) AddTo(t time.Time) time.Time {
	t = t.AddDate(d.Years, d.Months, d.Days+d.Weeks*7)
	t = t.Add(time.Duration(d.Hours) * time.Hour)
	t = t.Add(time.Duration(d.Minutes) * time.Minute)
	t = t.Add(time.Duration(d.Seconds) * time.Second)
	return t
}
