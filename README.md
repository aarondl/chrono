# chrono

A simple Go package for dividing time.Time into datetime, time, and date
such that it becomes difficult to use any of these types incorrectly to access
the excluded portions.

It is a goal of this package to expose an API surface that is very similar
to the time.Time type, and re-use things like time.Duration and time.Month.

In addition to Marshaling/Unmarshaling interfaces, it also supports SQL
interfaces directly, rather than relying on the driver's time.Time handling
behavior.

# Formatting

This package defaults to RFC3339 (ISO8601 compatible) inputs and outputs with
the exception of SQL handling in which case it attempts to be closer to the SQL standard dialect of ISO8601.

# Examples

```go
// Creating values
datetime := chrono.NewDateTime(2000, 1, 2, 3, 4, 5, 10, time.UTC)
date := chrono.NewDate(2000, 1, 2)
time := chrono.NewDateTime(3, 4, 5, 10, time.UTC)

// Unix timestamp interop
datetime = chrono.DateTimeFromUnix(1655610143, 0)
date = chrono.DateFromUnix(1655610143, 0)
time = chrono.TimeFromUnix(1655610143, 0)

// Creating values from strings (RFC3339), see 'Layout' variant functions for
// functionality from the std library for parsing custom formats.
datetime = chrono.DateTimeFromString("2000-01-02T03:04:05Z")
date = chrono.DateTimeFromString("2000-01-02")
time = chrono.DateTimeFromString("03:04:05Z")

// Conversion from/to time.Time
datetime = chrono.DateTimeFromStdTime(time.Now())
date = chrono.DateFromStdTime(time.Now())
time = chrono.TimeFromStdTime(time.Now())
datetime.ToStdTime()
date.ToStdTime()
time.ToStdTime()

// Conversions between DateTime/Date/Time
datetime = chrono.DateTimeFromStdTime(time.Now())
date = datetime.ToDate()
time = datetime.ToTime()

// Additional comparisons on top of the stdlib After/Before/Equal
datetime.AfterOrEqual(datetime)
datetime.BeforeOrEqual(datetime)
datetime.Between(datetime)
datetime.BetweenOrEqual(datetime)

// Formatting is all done with RFC3339 parts
datetime.String() // "2000-01-02T03:04:05Z"
date.String()     // "2000-01-02"
time.String()     // "03:04:05Z"
```
