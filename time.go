package chrono

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	timeLayout       = "15:04:05Z07:00"
	quotedTimeLayout = `"` + timeLayout + `"`
	// TimeSQLLayout is exported so you can change this for your project
	// but the default should be sufficient. It used microsecond precision
	// to align with postgresq/mysql.
	TimeSQLLayout = "15:04:05.999999-07"
)

// Time is mostly a pass-through wrapper for time.Time. This allows
// nicer interoperability with the Time and Date types as well as a couple
// additional utility methods.
type Time struct {
	t time.Time
}

// NewTime from all components
func NewTime(hour, min, sec, nsec int, loc *time.Location) Time {
	return Time{t: time.Date(0, 1, 1, hour, min, sec, nsec, time.UTC)}
}

// TimeFromNow creates a new date time from the current moment in time
// (local).
func TimeFromNow() Time {
	return Time{t: time.Now()}
}

// TimeFromString parses a date time (ISO8601/RFC3339 date-time) in the
// local location.
func TimeFromString(str string) (Time, error) {
	t, err := time.Parse(timeLayout, str)
	if err != nil {
		return Time{}, fmt.Errorf("failed to parse time (%s): %w", str, err)
	}

	return Time{t: t}, nil
}

// TimeFromStringLocation parses a date time (ISO8601/RFC3339 date-time) in
// the specified location.
func TimeFromStringLocation(str string, loc *time.Location) (Time, error) {
	t, err := time.ParseInLocation(timeLayout, str, loc)
	if err != nil {
		return Time{}, fmt.Errorf("failed to parse time (%s): %w", str, err)
	}

	return Time{t: t}, nil
}

// TimeFromString parses a time from a layout in the local location.
func TimeFromLayout(layout, str string) (Time, error) {
	t, err := time.Parse(layout, str)
	if err != nil {
		return Time{}, fmt.Errorf("failed to parse time (%s): %w", str, err)
	}

	return Time{t: t}, nil
}

// TimeFromStringLocation parses a time from a layout in the specified location.
func TimeFromLayoutLocation(layout, str string, loc *time.Location) (Time, error) {
	t, err := time.ParseInLocation(timeLayout, str, loc)
	if err != nil {
		return Time{}, fmt.Errorf("failed to parse time (%s): %w", str, err)
	}

	return Time{t: t}, nil
}

// TimeFromStdTime creates a time object discarding the stdlib time.Time's date
// information.
func TimeFromStdTime(t time.Time) Time {
	return Time{t: time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())}
}

// Unix returns the local Time corresponding to the given Unix time, discards
// the date information.
func TimeFromUnix(sec int64, nsec int64) Time {
	return TimeFromStdTime(time.Unix(sec, nsec).UTC())
}

// UnixMicro returns the local Time corresponding to the given Unix time in
// microseconds. Discards the date information.
func TimeFromUnixMicro(usec int64) Time {
	return TimeFromStdTime(time.UnixMicro(usec).UTC())
}

// UnixMicro returns the local Time corresponding to the given Unix time in
// milliseconds. Discards the date information.
func TimeFromUnixMilli(msec int64) Time {
	return TimeFromStdTime(time.UnixMilli(msec).UTC())
}

// ToStdTime returns the time as a time.Time
func (t Time) ToStdTime() time.Time {
	return time.Date(0, 1, 1, t.t.Hour(), t.t.Minute(), t.t.Second(), t.t.Nanosecond(), t.t.Location())
}

// Add returns the time t+d.
func (t Time) Add(dur time.Duration) Time {
	return TimeFromStdTime(t.t.Add(dur))
}

// After returns true if rhs is after d
func (t Time) After(rhs Time) bool {
	return t.t.After(rhs.t)
}

// AfterOrEqual returns true if rhs is equal to or after d
func (t Time) AfterOrEqual(rhs Time) bool {
	return t.t.After(rhs.t) || t.t.Equal(rhs.t)
}

// AppendFormat is like Format but appends the textual representation to b and
// returns the extended buffer. Due to this package using time.Time the layout
// string is not checked for date-like parts that could be leaked out but will
// be zero.
func (t Time) AppendFormat(b []byte, layout string) []byte {
	return t.t.AppendFormat(b, layout)
}

// Before returns true if rhs is before d
func (t Time) Before(rhs Time) bool {
	return t.t.Before(rhs.t)
}

// BeforeOrEqual returns true if rhs is before d
func (t Time) BeforeOrEqual(rhs Time) bool {
	return t.t.Before(rhs.t) || t.t.Equal(rhs.t)
}

// Between returns true if t is in the exclusive time range (start, end)
func (t Time) Between(start, end Time) bool {
	return t.t.After(start.t) && t.t.Before(end.t)
}

// BetweenOrEqual returns true if t is in the inclusive time range [start, end]
func (t Time) BetweenOrEqual(start, end Time) bool {
	return t.AfterOrEqual(start) && t.BeforeOrEqual(end)
}

// Equal returns true if rhs == d
func (t Time) Equal(rhs Time) bool {
	return t.t.Equal(rhs.t)
}

// GoString implements fmt.GoStringer
func (t Time) GoString() string {
	hr, min, sec := t.t.Clock()
	nsec := t.t.Nanosecond()
	return fmt.Sprintf("chrono.Time(%d, %d, %d, %d, %s)", hr, min, sec, nsec, t.t.Location())
}

// MarshalBinary implements the encoding.BinaryMarshaler interface. This is
// inefficient because it actually will use time.Time's entire MarshalBinary
// method which means that it will be much larger due to date information also
// being stored.
func (t Time) MarshalBinary() ([]byte, error) {
	return t.t.MarshalBinary()
}

// MarshalJSON implements json.Marshaller
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(t.t.Format(quotedTimeLayout)), nil
}

// MarshalText implements encoding.TextMarshaller
func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// String returns an ISO8601 Time, also an RFC3339 date-time
func (t Time) String() string {
	return t.t.Format(timeLayout)
}

// UnmarshalBinary
func (d *Time) UnmarshalBinary(data []byte) error {
	var t time.Time
	if err := t.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("failed to unmarshal Time (%q): %w", data, err)
	}
	d.t = t
	return nil
}

// UnmarshalJSON parses a quoted ISO8601 Time / RFC3339 full-time
func (d *Time) UnmarshalJSON(data []byte) error {
	t, err := time.Parse(quotedTimeLayout, string(data))
	if err != nil {
		return fmt.Errorf("failed to unmarshal time (%q): %w", data, err)
	}
	d.t = t
	return nil
}

// UnmarshalText parses a byte string with ISO8601 Time / RFC3339 full-time
func (d *Time) UnmarshalText(data []byte) error {
	t, err := time.Parse(timeLayout, string(data))
	if err != nil {
		return fmt.Errorf("failed to unmarshal time (%q): %w", data, err)
	}
	d.t = t
	return nil
}

// Clock returns the time components
func (t Time) Clock() (hour, min, sec int) {
	return t.t.Clock()
}

// Format using a layout string from time.Time. This can accidentally pull
// zero'd date information from the underlying time.Time so caution must be
// used.
func (t Time) Format(layout string) string {
	return t.t.Format(layout)
}

// Hour returns the hour
func (t Time) Hour() int {
	return t.t.Hour()
}

// In returns the Time in the specified location
func (t Time) In(loc *time.Location) Time {
	return Time{t: t.t.In(loc)}
}

// IsDST returns true if DST is active
func (t Time) IsDST() bool {
	return t.t.IsDST()
}

// IsZero returns true if the Date is the zero value.
func (t Time) IsZero() bool {
	return t.t.IsZero()
}

// Local returns the current date time in the local location
func (t Time) Local() Time {
	return Time{t: t.t.Local()}
}

// Location returns the Time's location
func (t Time) Location() *time.Location {
	return t.t.Location()
}

// Minute returns the minute of the hour
func (t Time) Minute() int {
	return t.t.Minute()
}

// Nanosecond returns the nanosecond offset
func (t Time) Nanosecond() int {
	return t.t.Nanosecond()
}

// Round to the duration unit specified
func (t Time) Round(dur time.Duration) Time {
	return Time{t: t.t.Round(dur)}
}

// Second returns the second of the minute
func (t Time) Second() int {
	return t.t.Second()
}

// Sub returns the duration between the two times
func (t Time) Sub(u Time) time.Duration {
	return t.t.Sub(u.t)
}

// Truncate to the duration unit specified
func (t Time) Truncate(dur time.Duration) Time {
	return Time{t: t.t.Truncate(dur)}
}

// UTC returns the date time in UTC
func (t Time) UTC() Time {
	return Time{t: t.t.UTC()}
}

func (t Time) Zone() (name string, offset int) {
	return t.t.Zone()
}

// Value implements driver.Valuer
func (t Time) Value() (driver.Value, error) {
	return t.t.Format(TimeSQLLayout), nil
}

// Scan implements sql.Scanner. SQL requires the use of ISO8601.
func (t *Time) Scan(value any) error {
	if value == nil {
		t.t = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case int64:
		// Assume this is a unix timestamp
		*t = TimeFromUnix(v, 0)
		return nil
	case float64:
		// Assume this is a unix timestamp in float
		*t = TimeFromUnix(int64(v), 0)
		return nil
	case string:
		newt, err := time.Parse(TimeSQLLayout, v)
		if err != nil {
			return fmt.Errorf("failed to scan time (%q): %w", v, err)
		}
		t.t = newt
		return nil
	case []byte:
		newt, err := time.Parse(TimeSQLLayout, string(v))
		if err != nil {
			return fmt.Errorf("failed to scan time (%q): %w", v, err)
		}
		t.t = newt
		return nil
	case time.Time:
		*t = TimeFromStdTime(v)
		return nil
	}

	return fmt.Errorf("failed to scan type '%T' into time", value)
}
