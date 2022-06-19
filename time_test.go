package chrono_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/aarondl/chrono"
)

func TestTimeConstructors(t *testing.T) {
	t.Parallel()

	ref := chrono.NewTime(3, 4, 5, 0, time.UTC)
	now := chrono.TimeFromNow()
	if ref.AfterOrEqual(now) {
		t.Error("should be after old time")
	}
	dt, err := chrono.TimeFromString("03:04:05Z")
	if err != nil {
		t.Error(err)
	}
	if !ref.Equal(dt) {
		t.Error("should be equal", dt, ref)
	}
	dt, err = chrono.TimeFromStringLocation("03:04:05Z", time.UTC)
	if err != nil {
		t.Error(err)
	}
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt, err = chrono.TimeFromLayout("15:04:05Z07:00", "03:04:05Z")
	if err != nil {
		t.Error(err)
	}
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt, err = chrono.TimeFromLayoutLocation("15:04:05Z07:00", "03:04:05Z", time.UTC)
	if err != nil {
		t.Error(err)
	}
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt = chrono.TimeFromUnix(946695845, 0)
	if !ref.Equal(dt) {
		t.Error("should be equal", dt, ref)
	}
	dt = chrono.TimeFromUnixMicro(946695845000000)
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt = chrono.TimeFromUnixMilli(946695845000)
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
}

func TestTimeConversions(t *testing.T) {
	t.Parallel()

	stdTime := time.Date(0, 1, 1, 3, 4, 5, 0, time.UTC)
	ref := chrono.NewTime(3, 4, 5, 0, time.UTC)

	dt := chrono.TimeFromStdTime(stdTime)
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	cmp := dt.ToStdTime()

	if !cmp.Equal(stdTime) {
		t.Error("should be equal")
	}
}

func TestTimeModifications(t *testing.T) {
	t.Parallel()

	ref := chrono.NewTime(3, 4, 30, 0, time.UTC)
	dt := ref.Add(time.Hour)
	if !dt.Equal(chrono.NewTime(4, 4, 30, 0, time.UTC)) {
		t.Error("should be equal", dt)
	}

	if ref.In(time.Local).Location() != time.Local {
		t.Error("should be in local")
	}
	if ref.Local().Location() != time.Local {
		t.Error("should be in local")
	}
	if ref.Local().UTC().Location() != time.UTC {
		t.Error("should be in UTC")
	}

	dt = ref.Round(time.Minute)
	if !dt.Equal(chrono.NewTime(3, 5, 0, 0, time.UTC)) {
		t.Error("should be equal", dt)
	}

	dt = ref.Truncate(time.Minute)
	if !dt.Equal(chrono.NewTime(3, 4, 0, 0, time.UTC)) {
		t.Error("should be equal", dt)
	}

	dur := ref.Sub(chrono.NewTime(3, 4, 0, 0, time.UTC))
	if dur != time.Second*30 {
		t.Error("wrong value")
	}
}

func TestTimeComparisons(t *testing.T) {
	t.Parallel()

	ref := chrono.NewTime(3, 4, 30, 0, time.UTC)

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
	if !chrono.TimeFromNow().After(ref) {
		t.Error("it should be after the ref time")
	}
	if !chrono.TimeFromNow().AfterOrEqual(ref) {
		t.Error("it should be after the ref time")
	}
	if ref.After(chrono.TimeFromNow()) {
		t.Error("ref should not be after now")
	}
	if ref.AfterOrEqual(chrono.TimeFromNow()) {
		t.Error("ref should not be after now")
	}

	// Before
	if !ref.Before(chrono.TimeFromNow()) {
		t.Error("it should be before the ref time")
	}
	if !ref.BeforeOrEqual(chrono.TimeFromNow()) {
		t.Error("it should be before the ref time")
	}
	if chrono.TimeFromNow().Before(ref) {
		t.Error("now should not be before the ref time")
	}
	if chrono.TimeFromNow().BeforeOrEqual(ref) {
		t.Error("now should not be before the ref time")
	}

	// Between
	before := chrono.NewTime(1, 0, 0, 0, time.UTC)
	after := chrono.NewTime(4, 0, 0, 0, time.UTC)
	if !ref.Between(before, after) {
		t.Error("it should be between")
	}
	if chrono.TimeFromNow().Between(before, after) {
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
	if chrono.TimeFromNow().BetweenOrEqual(before, after) {
		t.Error("now should not be between")
	}
	if !ref.BetweenOrEqual(ref, after) {
		t.Error("it should be between")
	}
	if !ref.BetweenOrEqual(before, ref) {
		t.Error("it should be between")
	}
}

func TestTimeFormatting(t *testing.T) {
	t.Parallel()

	ref := chrono.NewTime(3, 4, 30, 0, time.UTC)
	var b []byte
	if ob := ref.AppendFormat(b, "15:04:05Z07:00"); !bytes.Equal(ob, []byte("03:04:30Z")) {
		t.Error("bytes were wrong:", string(ob))
	}

	if ref.GoString() != "chrono.Time(3, 4, 30, 0, UTC)" {
		t.Error("string was wrong:", ref.GoString())
	}

	if ref.String() != "03:04:30Z" {
		t.Error("string was wrong:", ref.String())
	}

	if ref.Format("03:04:05Z07:00") != "03:04:30Z" {
		t.Error("string was wrong:", ref.String())
	}
}

func TestTimeGetters(t *testing.T) {
	t.Parallel()

	ref := chrono.NewTime(3, 4, 30, 10, time.UTC)

	if v := ref.Hour(); v != 3 {
		t.Error("value wrong:", v)
	}
	if v := ref.Minute(); v != 4 {
		t.Error("value wrong:", v)
	}
	if v := ref.Second(); v != 30 {
		t.Error("value wrong:", v)
	}
	if v := ref.Nanosecond(); v != 10 {
		t.Error("value wrong:", v)
	}
	if h, m, s := ref.Clock(); h != 3 || m != 4 || s != 30 {
		t.Error("value wrong:", h, m, s)
	}
	if ref.IsDST() {
		t.Error("no dst available for UTC")
	}
	if ref.IsZero() {
		t.Error("not zero")
	}
	if v := ref.Location(); v != time.UTC {
		t.Error("value wrong:", v)
	}
	if name, offset := ref.Zone(); name != "UTC" || offset != 0 {
		t.Error("value wrong:", name, offset)
	}
}

func TestTimeMarshalling(t *testing.T) {
	t.Parallel()

	ref := chrono.NewTime(3, 4, 30, 0, time.UTC)
	bin, err := ref.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	var unbin chrono.Time
	if err = unbin.UnmarshalBinary(bin); err != nil {
		t.Error(err)
	}
	if !unbin.Equal(ref) {
		t.Error("value was wrong")
	}

	js, err := ref.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	if string(js) != `"03:04:30Z"` {
		t.Error("value wrong", string(js))
	}
	var unjs chrono.Time
	if err = unjs.UnmarshalJSON(js); err != nil {
		t.Error(err)
	}
	if !unjs.Equal(ref) {
		t.Error("value was wrong", unjs, ref)
	}

	txt, err := ref.MarshalText()
	if err != nil {
		t.Error(err)
	}
	if string(txt) != `03:04:30Z` {
		t.Error("value wrong")
	}
	var untxt chrono.Time
	if err = untxt.UnmarshalText(txt); err != nil {
		t.Error(err)
	}
	if !untxt.Equal(ref) {
		t.Error("value was wrong")
	}
}

func TestTimeSQL(t *testing.T) {
	t.Parallel()

	ref := chrono.NewTime(3, 4, 5, 0, time.UTC)
	if v, err := ref.Value(); err != nil {
		t.Error(err)
	} else if v.(string) != "03:04:05+00" {
		t.Error("value was wrong", v)
	}

	var newt chrono.Time
	if err := newt.Scan("03:04:05+00"); err != nil {
		t.Error(err)
	}
	if !newt.Equal(ref) {
		t.Error("value was wrong")
	}

	newt = chrono.Time{}
	if err := newt.Scan([]byte("03:04:05+00")); err != nil {
		t.Error(err)
	}
	if !newt.Equal(ref) {
		t.Error("value was wrong")
	}

	newt = chrono.Time{}
	if err := newt.Scan(int64(946695845)); err != nil {
		t.Error(err)
	}
	if !newt.Equal(ref) {
		t.Error("value was wrong")
	}

	newt = chrono.Time{}
	if err := newt.Scan(float64(946695845)); err != nil {
		t.Error(err)
	}
	if !newt.Equal(ref) {
		t.Error("value was wrong")
	}

	newt = chrono.Time{}
	if err := newt.Scan(ref.ToStdTime()); err != nil {
		t.Error(err)
	}
	if !newt.Equal(ref) {
		t.Error("value was wrong")
	}
}
