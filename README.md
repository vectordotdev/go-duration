# go-duration

[![GoDoc](https://godoc.org/github.com/timberio/go-duration?status.svg)](http://godoc.org/github.com/timberio/go-duration)
[![Circle CI](https://circleci.com/gh/timberio/go-duration.svg?style=svg)](https://circleci.com/gh/timberio/go-duration)
[![Go Report Card](https://goreportcard.com/badge/github.com/timberio/go-duration)](https://goreportcard.com/report/github.com/timberio/go-duration)
[![coverage](https://gocover.io/_badge/github.com/timberio/go-duration?0 "coverage")](http://gocover.io/github.com/timberio/go-duration)

This library provides facilities for parsing and formatting durations as described by
[RFC3339](https://tools.ietf.org/html/rfc3339) which describes a specific encoding for
[ISO8601](https://tools.ietf.org/html/rfc3339#ref-ISO8601) dates, times, and time periods.

Example usage:

```go
d, _ := duration.ParseRFC3339("P3Y6M4DT12H30M5S")
fmt.Println(duration.AddTo(time.Now())
fmt.Println(duration.FormatRFC3339())
```

See [package documentation](http://godoc.org/github.com/timberio/go-duration) for usage and more examples.
