package chrono_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/aarondl/chrono"
)

func TestDateConstructors(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDate(2000, 1, 2)
	now := chrono.DateFromNow()
	if ref.AfterOrEqual(now) {
		t.Error("should be after old time")
	}
	dt, err := chrono.DateFromString("2000-01-02")
	if err != nil {
		t.Error(err)
	}
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt, err = chrono.DateFromLayout("2006-01-02", "2000-01-02")
	if err != nil {
		t.Error(err)
	}
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt = chrono.DateFromUnix(ref.Unix(), 0)
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt = chrono.DateFromUnixMicro(ref.UnixMicro())
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt = chrono.DateFromUnixMilli(ref.UnixMilli())
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
}

func TestDateConversions(t *testing.T) {
	t.Parallel()

	stdTime := time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)
	ref := chrono.NewDate(2000, 1, 2)

	dt := chrono.DateFromStdTime(stdTime)
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	cmp := dt.ToStdTime()

	if !cmp.Equal(stdTime) {
		t.Error("should be equal")
	}
}

func TestDateModifications(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDate(2000, 1, 2)

	dt := ref.AddDate(0, 0, 1)
	if !dt.Equal(chrono.NewDate(2000, 1, 3)) {
		t.Error("should be equal", dt)
	}
}

func TestDateComparisons(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDate(2000, 1, 2)

	// Equality
	if !ref.Equal(ref) {
		t.Error("should be equal")
	}
	if !ref.AfterOrEqual(ref) {
		t.Error("should be equal")
	}
	if !ref.BeforeOrEqual(ref) {
		t.Error("should be equal")
	}

	// After
	if !chrono.DateFromNow().After(ref) {
		t.Error("it should be after the ref date")
	}
	if !chrono.DateFromNow().AfterOrEqual(ref) {
		t.Error("it should be after the ref date")
	}
	if ref.After(chrono.DateFromNow()) {
		t.Error("ref should not be after now")
	}
	if ref.AfterOrEqual(chrono.DateFromNow()) {
		t.Error("ref should not be after now")
	}

	// Before
	if !ref.Before(chrono.DateFromNow()) {
		t.Error("it should be before the ref date")
	}
	if !ref.BeforeOrEqual(chrono.DateFromNow()) {
		t.Error("it should be before the ref date")
	}
	if chrono.DateFromNow().Before(ref) {
		t.Error("now should not be before the ref date")
	}
	if chrono.DateFromNow().BeforeOrEqual(ref) {
		t.Error("now should not be before the ref date")
	}

	// Between
	before := chrono.NewDate(2000, 1, 1)
	after := chrono.NewDate(2000, 1, 3)
	if !ref.Between(before, after) {
		t.Error("it should be between")
	}
	if chrono.DateFromNow().Between(before, after) {
		t.Error("now should not be between")
	}
	if ref.Between(ref, after) {
		t.Error("it should not be between because exclusive")
	}
	if ref.Between(before, ref) {
		t.Error("it should not be between")
	}
	if !ref.BetweenOrEqual(before, after) {
		t.Error("it should be between")
	}
	if chrono.DateFromNow().BetweenOrEqual(before, after) {
		t.Error("now should not be between")
	}
	if !ref.BetweenOrEqual(ref, after) {
		t.Error("it should be between")
	}
	if !ref.BetweenOrEqual(before, ref) {
		t.Error("it should be between")
	}
}

func TestDateFormatting(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDate(2000, 1, 2)
	var b []byte
	if ob := ref.AppendFormat(b, "2006-01-02"); !bytes.Equal(ob, []byte("2000-01-02")) {
		t.Error("bytes were wrong:", string(ob))
	}

	if ref.GoString() != "chrono.Date(2000, January, 2)" {
		t.Error("string was wrong:", ref.GoString())
	}

	if ref.String() != "2000-01-02" {
		t.Error("string was wrong:", ref.String())
	}

	if ref.Format("2006-01-02") != "2000-01-02" {
		t.Error("string was wrong:", ref.String())
	}
}

func TestDateGetters(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDate(2000, 1, 2)

	if y, m, d := ref.Date(); y != 2000 || m != 1 || d != 2 {
		t.Error("value wrong:", y, m, d)
	}
	if v := ref.Unix(); v != 946771200 {
		t.Error("value wrong:", v)
	}
	if v := ref.UnixMicro(); v != 946771200000000 {
		t.Error("value wrong:", v)
	}
	if v := ref.UnixMilli(); v != 946771200000 {
		t.Error("value wrong:", v)
	}
	if v := ref.UnixNano(); v != 946771200000000000 {
		t.Error("value wrong:", v)
	}
	if v := ref.Weekday(); v != time.Sunday {
		t.Error("value wrong:", v)
	}
	if v := ref.YearDay(); v != 2 {
		t.Error("value wrong:", v)
	}
	if v := ref.Year(); v != 2000 {
		t.Error("value wrong:", v)
	}
	if v := ref.Month(); v != 1 {
		t.Error("value wrong:", v)
	}
	if v := ref.Day(); v != 2 {
		t.Error("value wrong:", v)
	}
	if ref.IsZero() {
		t.Error("not zero")
	}
	// Awkward result, but this is an implementation detail of time.Time
	if year, week := ref.ISOWeek(); year != 1999 || week != 52 {
		t.Error("value wrong:", year, week)
	}
}

func TestDateMarshalling(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDate(2000, 1, 2)
	bin, err := ref.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	var unbin chrono.Date
	if err = unbin.UnmarshalBinary(bin); err != nil {
		t.Error(err)
	}
	if !unbin.Equal(ref) {
		t.Error("value was wrong", unbin, ref)
	}

	js, err := ref.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	if string(js) != `"2000-01-02"` {
		t.Error("value wrong")
	}
	var unjs chrono.Date
	if err = unjs.UnmarshalJSON(js); err != nil {
		t.Error(err)
	}
	if !unjs.Equal(ref) {
		t.Error("value was wrong")
	}

	txt, err := ref.MarshalText()
	if err != nil {
		t.Error(err)
	}
	if string(txt) != `2000-01-02` {
		t.Error("value wrong")
	}
	var untxt chrono.Date
	if err = untxt.UnmarshalText(txt); err != nil {
		t.Error(err)
	}
	if !untxt.Equal(ref) {
		t.Error("value was wrong")
	}
}

func TestDateSQL(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDate(2000, 1, 2)
	if v, err := ref.Value(); err != nil {
		t.Error(err)
	} else if v.(string) != "2000-01-02" {
		t.Error("value was wrong")
	}

	var date chrono.Date
	if err := date.Scan("2000-01-02"); err != nil {
		t.Error(err)
	}
	if !date.Equal(ref) {
		t.Error("value was wrong")
	}

	date = chrono.Date{}
	if err := date.Scan([]byte("2000-01-02")); err != nil {
		t.Error(err)
	}
	if !date.Equal(ref) {
		t.Error("value was wrong")
	}

	date = chrono.Date{}
	if err := date.Scan(int64(ref.Unix())); err != nil {
		t.Error(err)
	}
	if !date.Equal(ref) {
		t.Error("value was wrong")
	}

	date = chrono.Date{}
	if err := date.Scan(float64(ref.Unix())); err != nil {
		t.Error(err)
	}
	if !date.Equal(ref) {
		t.Error("value was wrong")
	}

	date = chrono.Date{}
	if err := date.Scan(ref.ToStdTime()); err != nil {
		t.Error(err)
	}
	if !date.Equal(ref) {
		t.Error("value was wrong")
	}
}

func TestAddMonthsNoOverflow(t *testing.T) {
	t.Parallel()

	t.Run("AddOneMonthToJanuary", func(t *testing.T) {
		ref := chrono.NewDate(2024, 1, 31)
		dt := ref.AddMonthsNoOverflow(1)
		if !dt.Equal(chrono.NewDate(2024, 2, 29)) {
			t.Error("should be equal", dt)
		}
	})

	t.Run("AddOneMonthToMarch", func(t *testing.T) {
		ref := chrono.NewDate(2024, 3, 31)
		dt := ref.AddMonthsNoOverflow(1)
		if !dt.Equal(chrono.NewDate(2024, 4, 30)) {
			t.Error("should be equal", dt)
		}
	})

	t.Run("AddOneMonthToMay", func(t *testing.T) {
		ref := chrono.NewDate(2024, 5, 31)
		dt := ref.AddMonthsNoOverflow(1)
		if !dt.Equal(chrono.NewDate(2024, 6, 30)) {
			t.Error("should be equal", dt)
		}
	})

	t.Run("AddOneMonthToAugust", func(t *testing.T) {
		ref := chrono.NewDate(2024, 8, 31)
		dt := ref.AddMonthsNoOverflow(1)
		if !dt.Equal(chrono.NewDate(2024, 9, 30)) {
			t.Error("should be equal", dt)
		}
	})

	t.Run("AddOneMonthToSeptember", func(t *testing.T) {
		ref := chrono.NewDate(2024, 10, 31)
		dt := ref.AddMonthsNoOverflow(1)
		if !dt.Equal(chrono.NewDate(2024, 11, 30)) {
			t.Error("should be equal", dt)
		}
	})
}

func TestSubtractMonthsNoOverflow(t *testing.T) {
	t.Parallel()

	t.Run("SubtractOneMonthFromMarch", func(t *testing.T) {
		ref := chrono.NewDate(2024, 3, 31)
		dt := ref.AddMonthsNoOverflow(-1)
		if !dt.Equal(chrono.NewDate(2024, 2, 29)) {
			t.Error("should be equal", dt)
		}
	})

	t.Run("SubtractOneMonthFromMai", func(t *testing.T) {
		ref := chrono.NewDate(2024, 5, 31)
		dt := ref.AddMonthsNoOverflow(-1)
		if !dt.Equal(chrono.NewDate(2024, 4, 30)) {
			t.Error("should be equal", dt)
		}
	})

	t.Run("SubtractOneMonthFromJuly", func(t *testing.T) {
		ref := chrono.NewDate(2024, 7, 31)
		dt := ref.AddMonthsNoOverflow(-1)
		if !dt.Equal(chrono.NewDate(2024, 6, 30)) {
			t.Error("should be equal", dt)
		}
	})

	t.Run("SubtractOneMonthFromOctober", func(t *testing.T) {
		ref := chrono.NewDate(2024, 10, 31)
		dt := ref.AddMonthsNoOverflow(-1)
		if !dt.Equal(chrono.NewDate(2024, 9, 30)) {
			t.Error("should be equal", dt)
		}
	})

	t.Run("SubtractOneMonthFromDecember", func(t *testing.T) {
		ref := chrono.NewDate(2024, 12, 31)
		dt := ref.AddMonthsNoOverflow(-1)
		if !dt.Equal(chrono.NewDate(2024, 11, 30)) {
			t.Error("should be equal", dt)
		}
	})
}
