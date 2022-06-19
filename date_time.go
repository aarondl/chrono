package chrono

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	dateTimeSQLLayout = "2006-01-02 15:04:05-07"
)

// DateTime is mostly a pass-through wrapper for time.Time. This allows
// nicer interoperability with the Time and Date types as well as a couple
// additional utility methods.
type DateTime struct {
	t time.Time
}

// NewDateTime from all components
func NewDateTime(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) DateTime {
	return DateTime{t: time.Date(year, month, day, hour, min, sec, nsec, loc)}
}

// DateTimeFromNow creates a new date time from the current moment in time
// (local).
func DateTimeFromNow() DateTime {
	return DateTime{t: time.Now()}
}

// DateTimeFromString parses a date time (ISO8601/RFC3339 date-time) in the
// local location.
func DateTimeFromString(str string) (DateTime, error) {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return DateTime{}, fmt.Errorf("failed to parse datetime (%s): %w", str, err)
	}

	return DateTime{t: t}, nil
}

// DateTimeFromStringLocation parses a date time (ISO8601/RFC3339 date-time) in
// the specified location.
func DateTimeFromStringLocation(str string, loc *time.Location) (DateTime, error) {
	t, err := time.ParseInLocation(time.RFC3339, str, loc)
	if err != nil {
		return DateTime{}, fmt.Errorf("failed to parse datetime (%s): %w", str, err)
	}

	return DateTime{t: t}, nil
}

// DateTimeFromString parses a date time by layout in the local location.
func DateTimeFromLayout(layout, str string) (DateTime, error) {
	t, err := time.Parse(layout, str)
	if err != nil {
		return DateTime{}, fmt.Errorf("failed to parse datetime (%s): %w", str, err)
	}

	return DateTime{t: t}, nil
}

// DateTimeFromStringLocation parses a date time by layout in the specified
// location.
func DateTimeFromLayoutLocation(layout, str string, loc *time.Location) (DateTime, error) {
	t, err := time.ParseInLocation(layout, str, loc)
	if err != nil {
		return DateTime{}, fmt.Errorf("failed to parse datetime (%s): %w", str, err)
	}

	return DateTime{t: t}, nil
}

// Unix returns the local Time corresponding to the given Unix time
func DateTimeFromUnix(sec int64, nsec int64) DateTime {
	return DateTime{t: time.Unix(sec, nsec)}
}

// UnixMicro returns the local Time corresponding to the given Unix time in
// microseconds
func DateTimeFromUnixMicro(usec int64) DateTime {
	return DateTime{t: time.UnixMicro(usec)}
}

// UnixMicro returns the local Time corresponding to the given Unix time in
// milliseconds
func DateTimeFromUnixMilli(msec int64) DateTime {
	return DateTime{t: time.UnixMilli(msec)}
}

// DateTimeFromStdTime converts a time.Time into a datetime
func DateTimeFromStdTime(t time.Time) DateTime {
	return DateTime{t: t}
}

// ToStdTime returns the same moment in time as a time.Time
func (d DateTime) ToStdTime() time.Time {
	return d.t
}

// Add returns the time t+d.
func (d DateTime) Add(dur time.Duration) DateTime {
	return DateTime{t: d.t.Add(dur)}
}

// AddDate to t and return
func (d DateTime) AddDate(years int, months int, days int) DateTime {
	return DateTime{t: d.t.AddDate(years, months, days)}
}

// After returns true if rhs is after d
func (d DateTime) After(rhs DateTime) bool {
	return d.t.After(rhs.t)
}

// AfterOrEqual returns true if rhs is equal to or after d
func (d DateTime) AfterOrEqual(rhs DateTime) bool {
	return d.t.After(rhs.t) || d.t.Equal(rhs.t)
}

// AppendFormat passes through to the underlying time.Time but.
func (d DateTime) AppendFormat(b []byte, layout string) []byte {
	return d.t.AppendFormat(b, layout)
}

// Before returns true if rhs is before d
func (d DateTime) Before(rhs DateTime) bool {
	return d.t.Before(rhs.t)
}

// BeforeOrEqual returns true if rhs is before d
func (d DateTime) BeforeOrEqual(rhs DateTime) bool {
	return d.t.Before(rhs.t) || d.t.Equal(rhs.t)
}

// Date returns the DateTime's components
func (d DateTime) Date() (year int, month time.Month, day int) {
	return d.t.Date()
}

// Day returns the day of the month
func (d DateTime) Day() int {
	return d.t.Day()
}

// Equal returns true if rhs == d
func (d DateTime) Equal(rhs DateTime) bool {
	return d.t.Equal(rhs.t)
}

// GoString implements fmt.GoStringer
func (d DateTime) GoString() string {
	y, m, day := d.t.Date()
	hr, min, sec := d.t.Clock()
	nsec := d.t.Nanosecond()
	return fmt.Sprintf("chrono.DateTime(%d, %s, %d, %d, %d, %d, %d, %s)", y, m, day, hr, min, sec, nsec, d.t.Location())
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d DateTime) MarshalBinary() ([]byte, error) {
	return d.t.MarshalBinary()
}

// MarshalJSON implements json.Marshaller
func (d DateTime) MarshalJSON() ([]byte, error) {
	return d.t.MarshalJSON()
}

// MarshalText implements encoding.TextMarshaller
func (d DateTime) MarshalText() ([]byte, error) {
	return d.t.MarshalText()
}

// Month returns the month
func (d DateTime) Month() time.Month {
	return d.t.Month()
}

// String returns an ISO8601 DateTime, also an RFC3339 date-time
func (d DateTime) String() string {
	return d.t.Format(time.RFC3339)
}

// Unix timestamp
func (d DateTime) Unix() int64 {
	return d.t.Unix()
}

// UnixMicro returns a unix timestamp in microseconds
func (d DateTime) UnixMicro() int64 {
	return d.t.UnixMicro()
}

// UnixMilli returns a unix timestamp in milliseconds
func (d DateTime) UnixMilli() int64 {
	return d.t.UnixMilli()
}

// UnixNano returns a unix timestamp in nanoseconds
func (d DateTime) UnixNano() int64 {
	return d.t.UnixNano()
}

// UnmarshalBinary
func (d *DateTime) UnmarshalBinary(data []byte) error {
	var t time.Time
	if err := t.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("failed to unmarshal DateTime (%q): %w", data, err)
	}
	d.t = t
	return nil
}

// UnmarshalJSON parses a quoted ISO8601 DateTime / RFC3339 full-DateTime
func (d *DateTime) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := t.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("failed to unmarshal DateTime (%q): %w", data, err)
	}
	d.t = t
	return nil
}

// UnmarshalText parses a byte string with ISO8601 DateTime / RFC3339 full-DateTime
func (d *DateTime) UnmarshalText(data []byte) error {
	var t time.Time
	if err := t.UnmarshalText(data); err != nil {
		return fmt.Errorf("failed to unmarshal DateTime (%q): %w", data, err)
	}
	d.t = t
	return nil
}

// Weekday returns the day of the week
func (d DateTime) Weekday() time.Weekday {
	return d.t.Weekday()
}

// Year returns the year
func (d DateTime) Year() int {
	return d.t.Year()
}

// YearDay returns the day of the year
func (d DateTime) YearDay() int {
	return d.t.YearDay()
}

// Clock returns the time components
func (d DateTime) Clock() (hour, min, sec int) {
	return d.t.Clock()
}

// Format using a layout string from, same as time.Time
func (d DateTime) Format(layout string) string {
	return d.t.Format(layout)
}

// GobDecode passthrough
func (d *DateTime) GobDecode(data []byte) error {
	return d.t.GobDecode(data)
}

// GobEncode passthrough
func (d DateTime) GobEncode() ([]byte, error) {
	return d.t.GobEncode()
}

// Hour returns the hour
func (d DateTime) Hour() int {
	return d.t.Hour()
}

// ISOWeek returns the iso week
func (d DateTime) ISOWeek() (year, week int) {
	return d.t.ISOWeek()
}

// In returns the DateTime in the specified location
func (d DateTime) In(loc *time.Location) DateTime {
	return DateTime{t: d.t.In(loc)}
}

// IsDST returns true if DST is active
func (d DateTime) IsDST() bool {
	return d.t.IsDST()
}

// IsZero returns true if the Date is the zero value.
func (d DateTime) IsZero() bool {
	return d.t.IsZero()
}

// Local returns the current date time in the local location
func (d DateTime) Local() DateTime {
	return DateTime{t: d.t.Local()}
}

// Location returns the DateTime's location
func (d DateTime) Location() *time.Location {
	return d.t.Location()
}

// Minute returns the minute of the hour
func (d DateTime) Minute() int {
	return d.t.Minute()
}

// Nanosecond returns the nanosecond offset
func (d DateTime) Nanosecond() int {
	return d.t.Nanosecond()
}

// Round to the duration unit specified
func (d DateTime) Round(dur time.Duration) DateTime {
	return DateTime{t: d.t.Round(dur)}
}

// Second returns the second of the minute
func (d DateTime) Second() int {
	return d.t.Second()
}

// Sub returns the duration between the two times
func (d DateTime) Sub(u DateTime) time.Duration {
	return d.t.Sub(u.t)
}

// Truncate to the duration unit specified
func (d DateTime) Truncate(dur time.Duration) DateTime {
	return DateTime{t: d.t.Truncate(dur)}
}

// UTC returns the date time in UTC
func (d DateTime) UTC() DateTime {
	return DateTime{t: d.t.UTC()}
}

func (d DateTime) Zone() (name string, offset int) {
	return d.t.Zone()
}

// Value implements driver.Valuer. SQL requires the use of ISO8601.
func (d DateTime) Value() (driver.Value, error) {
	return d.t.Format(dateTimeSQLLayout), nil
}

// Scan implements sql.Scanner. SQL requires the use of ISO8601.
func (d *DateTime) Scan(value any) error {
	switch v := value.(type) {
	case int64:
		// Assume this is a unix timestamp
		d.t = time.Unix(v, 0).UTC()
		return nil
	case float64:
		// Assume this is a unix timestamp in float
		d.t = time.Unix(int64(v), 0).UTC()
		return nil
	case string:
		t, err := time.Parse(dateTimeSQLLayout, v)
		if err != nil {
			return fmt.Errorf("failed to scan datetime (%q): %w", v, err)
		}
		d.t = t
		return nil
	case []byte:
		t, err := time.Parse(dateTimeSQLLayout, string(v))
		if err != nil {
			return fmt.Errorf("failed to scan datetime (%q): %w", v, err)
		}
		d.t = t
		return nil
	}

	return fmt.Errorf("failed to scan type '%T' into datetime", value)
}
