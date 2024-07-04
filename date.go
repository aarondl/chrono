package chrono

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const (
	dateLayout       = "2006-01-02"
	quotedDateLayout = `"` + dateLayout + `"`
)

// Date type, based on time.Time.
type Date struct {
	t time.Time
}

// NewDate constructs a new date object from its components
func NewDate(year int, month time.Month, day int) Date {
	return Date{t: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// DateFromNow returns a new date using the current date. It uses time.Now()
// as a reference date, discarding time information.
func DateFromNow() Date {
	// Careful to use local time else we might end up changing dates
	// which would be unexpected.
	return DateFromStdTime(time.Now())
}

// DateFromString parses a Date from RFC3339 full-date
func DateFromString(str string) (Date, error) {
	t, err := time.ParseInLocation(dateLayout, str, time.UTC)
	if err != nil {
		return Date{}, fmt.Errorf("failed to parse date: %w", err)
	}

	return DateFromStdTime(t), nil
}

// DateFromLayout parses a Date from layout
func DateFromLayout(layout, str string) (Date, error) {
	t, err := time.ParseInLocation(layout, str, time.UTC)
	if err != nil {
		return Date{}, fmt.Errorf("failed to parse date: %w", err)
	}

	return DateFromStdTime(t), nil
}

// FromTime converts from the stdlib time.Time type, discarding time information
func DateFromStdTime(t time.Time) Date {
	return Date{t: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)}
}

// DateFromUnix converts a unix timestamp in seconds into a date.
func DateFromUnix(sec int64, nsec int64) Date {
	return DateFromStdTime(time.Unix(sec, nsec).UTC())
}

// DateFromUnixMicro converts a unix timestamp in microseconds into a date.
func DateFromUnixMicro(usec int64) Date {
	return DateFromStdTime(time.UnixMicro(usec).UTC())
}

// DateFromUnixMilli converts a unix timestamp in milliseconds into a date.
func DateFromUnixMilli(msec int64) Date {
	return DateFromStdTime(time.UnixMilli(msec).UTC())
}

// ToStdTime returns a time.Time with the time component zero'd out in UTC
// location.
func (d Date) ToStdTime() time.Time {
	// ensure we make a new one
	return time.Date(d.t.Year(), d.t.Month(), d.t.Day(), 0, 0, 0, 0, time.UTC)
}

// AddDate to the current date
func (d Date) AddDate(years int, months int, days int) Date {
	return DateFromStdTime(d.t.AddDate(years, months, days))
}

// AddMonthsNoOverflow adds a month to the current time, not overflowing in case the
// destination month has less days than the current one.
// Positive value travels forward while negative value travels into the past.
func (d Date) AddMonthsNoOverflow(m int) Date {
	addedDate := d.AddDate(0, m, 0)
	if d.Day() != addedDate.Day() {
		return addedDate.PreviousMonthLastDay()
	}

	return addedDate
}

// PreviousMonthLastDay returns the last day of the previous month.
func (d Date) PreviousMonthLastDay() Date {
	year, month, _ := d.Date()
	return NewDate(year, month, 0) // 0 makes it wrap to last month
}

// After returns true if rhs is after d
func (d Date) After(rhs Date) bool {
	return d.t.After(rhs.t)
}

// AfterOrEqual returns true if rhs is equal to or after d
func (d Date) AfterOrEqual(rhs Date) bool {
	return d.t.After(rhs.t) || d.t.Equal(rhs.t)
}

// AppendFormat is like Format but appends the textual representation to b and
// returns the extended buffer. Due to this package using time.Time the layout
// string is not checked for time-like parts that could be leaked out but will
// be zero.
func (d Date) AppendFormat(b []byte, layout string) []byte {
	return d.t.AppendFormat(b, layout)
}

// Before returns true if rhs is before d
func (d Date) Before(rhs Date) bool {
	return d.t.Before(rhs.t)
}

// BeforeOrEqual returns true if rhs is before d
func (d Date) BeforeOrEqual(rhs Date) bool {
	return d.t.Before(rhs.t) || d.t.Equal(rhs.t)
}

// Between returns true if d is in the exclusive time range (start, end)
func (d Date) Between(start, end Date) bool {
	return d.t.After(start.t) && d.t.Before(end.t)
}

// BetweenOrEqual returns true if d is in the inclusive time range [start, end]
func (d Date) BetweenOrEqual(start, end Date) bool {
	return d.AfterOrEqual(start) && d.BeforeOrEqual(end)
}

// Date returns the date's components
func (d Date) Date() (year int, month time.Month, day int) {
	return d.t.Date()
}

// Day returns the day of the month
func (d Date) Day() int {
	return d.t.Day()
}

// Equal returns true if rhs == d
func (d Date) Equal(rhs Date) bool {
	return d.t.Equal(rhs.t)
}

// Format using a layout string from time.Time. This can accidentally pull
// zero'd time information from the underlying time.Time so caution must be
// used.
func (d Date) Format(layout string) string {
	return d.t.Format(layout)
}

// GoString implements fmt.GoStringer
func (d Date) GoString() string {
	y, m, day := d.t.Date()
	return fmt.Sprintf("chrono.Date(%d, %s, %d)", y, m, day)
}

// IsZero returns true if the Date is the zero value.
func (d Date) IsZero() bool {
	return d.t.IsZero()
}

// MarshalBinary implements the encoding.BinaryMarshaler interface. Is always
// a width of 32 bits (4 bytes).
func (d Date) MarshalBinary() ([]byte, error) {
	var out uint32
	y, m, day := d.t.Date()
	// Year = 14 bits
	// Month = 4 bits
	// Day = 5 bits
	out |= uint32(y)
	out |= uint32(m) << 14
	out |= uint32(day) << (14 + 4)
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, out)
	return buf, nil
}

// MarshalJSON implements json.Marshaller
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d)), nil
}

// MarshalText implements encoding.TextMarshaller
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// Month returns the month
func (d Date) Month() time.Month {
	return d.t.Month()
}

// String returns an ISO8601 Date, also an RFC3339 full-date
func (d Date) String() string {
	return d.t.Format(dateLayout)
}

// Unix timestamp
func (d Date) Unix() int64 {
	return d.t.Unix()
}

// UnixMicro returns a unix timestamp in microseconds
func (d Date) UnixMicro() int64 {
	return d.t.UnixMicro()
}

// UnixMilli returns a unix timestamp in milliseconds
func (d Date) UnixMilli() int64 {
	return d.t.UnixMilli()
}

// UnixNano returns a unix timestamp in nanoseconds
func (d Date) UnixNano() int64 {
	return d.t.UnixNano()
}

// UnmarshalBinary
func (d *Date) UnmarshalBinary(data []byte) error {
	if len(data) != 4 {
		return errors.New("failed to unmarshal date, incorrect number of bytes")
	}
	in := binary.LittleEndian.Uint32(data)
	y, m, day := in&0b11_1111_1111_1111, (in>>14)&0b1111, (in>>(14+4))&0b1_1111
	*d = NewDate(int(y), time.Month(m), int(day))
	return nil
}

// UnmarshalJSON parses a quoted ISO8601 date / RFC3339 full-date
func (d *Date) UnmarshalJSON(data []byte) error {
	t, err := time.Parse(quotedDateLayout, string(data))
	if err != nil {
		return fmt.Errorf("failed to unmarshal date (%q): %w", data, err)
	}
	*d = DateFromStdTime(t)
	return nil
}

// UnmarshalText parses a byte string with ISO8601 date / RFC3339 full-date
func (d *Date) UnmarshalText(data []byte) error {
	t, err := time.Parse(dateLayout, string(data))
	if err != nil {
		return fmt.Errorf("failed to unmarshal date (%q): %w", data, err)
	}
	*d = DateFromStdTime(t)
	return nil
}

// Weekday returns the day of the week
func (d Date) Weekday() time.Weekday {
	return d.t.Weekday()
}

// Year returns the year
func (d Date) Year() int {
	return d.t.Year()
}

// YearDay returns the day of the year
func (d Date) YearDay() int {
	return d.t.YearDay()
}

// ISOWeek returns the ISO 8601 year and week numbers.
func (d Date) ISOWeek() (year, week int) {
	return d.t.ISOWeek()
}

// Value implements driver.Valuer. SQL requires the use of ISO8601.
func (d Date) Value() (driver.Value, error) {
	return d.t.Format(dateLayout), nil
}

// Scan implements sql.Scanner. SQL requires the use of ISO8601.
func (d *Date) Scan(value any) error {
	if value == nil {
		d.t = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case int64:
		// Assume this is a unix timestamp
		*d = DateFromUnix(v, 0)
		return nil
	case float64:
		// Assume this is a unix timestamp in float
		*d = DateFromUnix(int64(v), 0)
		return nil
	case string:
		t, err := time.Parse(dateLayout, v)
		if err != nil {
			return fmt.Errorf("failed to scan date (%q): %w", v, err)
		}
		d.t = t
		return nil
	case []byte:
		t, err := time.Parse(dateLayout, string(v))
		if err != nil {
			return fmt.Errorf("failed to scan date (%q): %w", v, err)
		}
		d.t = t
		return nil
	case time.Time:
		*d = DateFromStdTime(v)
		return nil
	}

	return fmt.Errorf("failed to scan type '%T' into date", value)
}
