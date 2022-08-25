package chrono_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/aarondl/chrono"
)

func TestDateTimeConstructors(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDateTime(2000, 1, 2, 3, 4, 5, 0, time.UTC)
	now := chrono.DateTimeFromNow()
	if ref.AfterOrEqual(now) {
		t.Error("should be after old time")
	}
	dt, err := chrono.DateTimeFromString("2000-01-02T03:04:05Z")
	if err != nil {
		t.Error(err)
	}
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt, err = chrono.DateTimeFromStringLocation("2000-01-02T03:04:05Z", time.UTC)
	if err != nil {
		t.Error(err)
	}
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt, err = chrono.DateTimeFromLayout(time.RFC3339, "2000-01-02T03:04:05Z")
	if err != nil {
		t.Error(err)
	}
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt, err = chrono.DateTimeFromLayoutLocation(time.RFC3339, "2000-01-02T03:04:05Z", time.UTC)
	if err != nil {
		t.Error(err)
	}
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt = chrono.DateTimeFromUnix(ref.Unix(), 0)
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt = chrono.DateTimeFromUnixMicro(ref.UnixMicro())
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	dt = chrono.DateTimeFromUnixMilli(ref.UnixMilli())
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
}

func TestDateTimeConversions(t *testing.T) {
	t.Parallel()

	stdTime := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
	ref := chrono.NewDateTime(2000, 1, 2, 3, 4, 5, 0, time.UTC)

	dt := chrono.DateTimeFromStdTime(stdTime)
	if !ref.Equal(dt) {
		t.Error("should be equal")
	}
	cmp := dt.ToStdTime()

	if !cmp.Equal(stdTime) {
		t.Error("should be equal")
	}

	cnvTime := dt.ToTime()
	cnvDate := dt.ToDate()

	ya, ma, da := cnvDate.Date()
	if ya != stdTime.Year() {
		t.Error("year was wrong:", ya)
	}
	if ma != stdTime.Month() {
		t.Error("month was wrong:", ma)
	}
	if da != stdTime.Day() {
		t.Error("day was wrong:", da)
	}

	if cnvTime.Hour() != stdTime.Hour() {
		t.Error("hour is wrong")
	}
	if cnvTime.Minute() != stdTime.Minute() {
		t.Error("min is wrong")
	}
	if cnvTime.Second() != stdTime.Second() {
		t.Error("second is wrong")
	}
	if cnvTime.Nanosecond() != stdTime.Nanosecond() {
		t.Error("ns is wrong")
	}
	if cnvTime.Location() != stdTime.Location() {
		t.Error("location is wrong")
	}
}

func TestDateTimeModifications(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDateTime(2000, 1, 2, 3, 4, 30, 0, time.UTC)
	dt := ref.Add(time.Hour)
	if !dt.Equal(chrono.NewDateTime(2000, 1, 2, 4, 4, 30, 0, time.UTC)) {
		t.Error("should be equal", dt)
	}

	dt = ref.AddDate(0, 0, 1)
	if !dt.Equal(chrono.NewDateTime(2000, 1, 3, 3, 4, 30, 0, time.UTC)) {
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
	if !dt.Equal(chrono.NewDateTime(2000, 1, 2, 3, 5, 0, 0, time.UTC)) {
		t.Error("should be equal", dt)
	}

	dt = ref.Truncate(time.Minute)
	if !dt.Equal(chrono.NewDateTime(2000, 1, 2, 3, 4, 0, 0, time.UTC)) {
		t.Error("should be equal", dt)
	}

	dur := ref.Sub(chrono.NewDateTime(2000, 1, 2, 3, 4, 0, 0, time.UTC))
	if dur != time.Second*30 {
		t.Error("wrong value")
	}
}

func TestDateTimeComparisons(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDateTime(2000, 1, 2, 3, 4, 30, 0, time.UTC)

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
	if !chrono.DateTimeFromNow().After(ref) {
		t.Error("it should be after the ref date")
	}
	if !chrono.DateTimeFromNow().AfterOrEqual(ref) {
		t.Error("it should be after the ref date")
	}
	if ref.After(chrono.DateTimeFromNow()) {
		t.Error("ref should not be after now")
	}
	if ref.AfterOrEqual(chrono.DateTimeFromNow()) {
		t.Error("ref should not be after now")
	}

	// Before
	if !ref.Before(chrono.DateTimeFromNow()) {
		t.Error("it should be before the ref date")
	}
	if !ref.BeforeOrEqual(chrono.DateTimeFromNow()) {
		t.Error("it should be before the ref date")
	}
	if chrono.DateTimeFromNow().Before(ref) {
		t.Error("now should not be before the ref date")
	}
	if chrono.DateTimeFromNow().BeforeOrEqual(ref) {
		t.Error("now should not be before the ref date")
	}

	// Between
	before := chrono.NewDateTime(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	after := chrono.NewDateTime(2000, 1, 3, 0, 0, 0, 0, time.UTC)
	if !ref.Between(before, after) {
		t.Error("it should be between")
	}
	if chrono.DateTimeFromNow().Between(before, after) {
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
	if chrono.DateTimeFromNow().BetweenOrEqual(before, after) {
		t.Error("now should not be between")
	}
	if !ref.BetweenOrEqual(ref, after) {
		t.Error("it should be between")
	}
	if !ref.BetweenOrEqual(before, ref) {
		t.Error("it should be between")
	}
}

func TestDateTimeFormatting(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDateTime(2000, 1, 2, 3, 4, 30, 0, time.UTC)
	var b []byte
	if ob := ref.AppendFormat(b, time.RFC3339); !bytes.Equal(ob, []byte("2000-01-02T03:04:30Z")) {
		t.Error("bytes were wrong:", string(ob))
	}

	if ref.GoString() != "chrono.DateTime(2000, January, 2, 3, 4, 30, 0, UTC)" {
		t.Error("string was wrong:", ref.GoString())
	}

	if ref.String() != "2000-01-02T03:04:30Z" {
		t.Error("string was wrong:", ref.String())
	}

	if ref.Format(time.RFC3339) != "2000-01-02T03:04:30Z" {
		t.Error("string was wrong:", ref.String())
	}
}

func TestDateTimeGetters(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDateTime(2000, 1, 2, 3, 4, 30, 10, time.UTC)

	if y, m, d := ref.Date(); y != 2000 || m != 1 || d != 2 {
		t.Error("value wrong:", y, m, d)
	}
	if v := ref.Unix(); v != 946782270 {
		t.Error("value wrong:", v)
	}
	if v := ref.UnixMicro(); v != 946782270000000 {
		t.Error("value wrong:", v)
	}
	if v := ref.UnixMilli(); v != 946782270000 {
		t.Error("value wrong:", v)
	}
	if v := ref.UnixNano(); v != 946782270000000010 {
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
	// Awkward result, but this is an implementation detail of time.Time
	if year, week := ref.ISOWeek(); year != 1999 || week != 52 {
		t.Error("value wrong:", year, week)
	}
}

func TestDateTimeMarshalling(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDateTime(2000, 1, 2, 3, 4, 30, 10, time.UTC)
	bin, err := ref.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	var unbin chrono.DateTime
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
	if string(js) != `"2000-01-02T03:04:30.00000001Z"` {
		t.Error("value wrong")
	}
	var unjs chrono.DateTime
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
	if string(txt) != `2000-01-02T03:04:30.00000001Z` {
		t.Error("value wrong")
	}
	var untxt chrono.DateTime
	if err = untxt.UnmarshalText(txt); err != nil {
		t.Error(err)
	}
	if !untxt.Equal(ref) {
		t.Error("value was wrong")
	}

	gob, err := ref.GobEncode()
	if err != nil {
		t.Error(err)
	}
	var ungob chrono.DateTime
	if err = ungob.GobDecode(gob); err != nil {
		t.Error(err)
	}
	if !ungob.Equal(ref) {
		t.Error("value was wrong")
	}
}

func TestDateTimeSQL(t *testing.T) {
	t.Parallel()

	ref := chrono.NewDateTime(2000, 1, 2, 3, 4, 5, 0, time.UTC)
	if v, err := ref.Value(); err != nil {
		t.Error(err)
	} else if v.(string) != "2000-01-02 03:04:05+00" {
		t.Error("value was wrong", v)
	}

	var datetime chrono.DateTime
	if err := datetime.Scan("2000-01-02 03:04:05+00"); err != nil {
		t.Error(err)
	}
	if !datetime.Equal(ref) {
		t.Error("value was wrong")
	}

	datetime = chrono.DateTime{}
	if err := datetime.Scan([]byte("2000-01-02 03:04:05+00")); err != nil {
		t.Error(err)
	}
	if !datetime.Equal(ref) {
		t.Error("value was wrong")
	}

	datetime = chrono.DateTime{}
	if err := datetime.Scan(int64(ref.Unix())); err != nil {
		t.Error(err)
	}
	if !datetime.Equal(ref) {
		t.Error("value was wrong")
	}

	datetime = chrono.DateTime{}
	if err := datetime.Scan(float64(ref.Unix())); err != nil {
		t.Error(err)
	}
	if !datetime.Equal(ref) {
		t.Error("value was wrong")
	}

	datetime = chrono.DateTime{}
	if err := datetime.Scan(ref.ToStdTime()); err != nil {
		t.Error(err)
	}
	if !datetime.Equal(ref) {
		t.Error("value was wrong")
	}
}
