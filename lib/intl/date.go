package intl

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/awoodbeck/strftime"
	"github.com/getevo/evo-ng/lib/generic"
	"github.com/getevo/monday"
	"golang.org/x/text/language"
	"strconv"
	"strings"
	"time"
)

type Time time.Time

// Date parse any kind of date format
//  @accept  year , month , day
//  @accept  year , month , day , hour , minute , second
//  @accept  year , month , day , hour , minute , second , nanosecond
//  @accept  year , month , day , hour , minute , second , nanosecond , location
//  @accept date string
//  @accept unix timestamp
//  @accept time.Time
//  @accept intl.Time
//  @param params... interface{}
//  @return Time
func Date(params ...interface{}) Time {
	length := len(params)
	if length == 1 {
		var d, _ = TryParseTime(params[0])
		return d
	} else if length == 3 {
		y := generic.Parse(params[0]).Int()
		var m = 1
		if v, ok := params[1].(int); ok {
			m = v
		} else if v, ok := params[1].(string); ok {
			m = generic.Parse(v).Int()
		}
		d := generic.Parse(params[2]).Int()
		return Time(time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC))
	} else if length == 6 {
		y := generic.Parse(params[0]).Int()
		var m = 1
		if v, ok := params[1].(int); ok {
			m = v
		} else if v, ok := params[1].(string); ok {
			m = generic.Parse(v).Int()
		}
		d := generic.Parse(params[2]).Int()
		h := generic.Parse(params[3]).Int()
		i := generic.Parse(params[4]).Int()
		s := generic.Parse(params[5]).Int()
		return Time(time.Date(y, time.Month(m), d, h, i, s, 0, time.UTC))
	} else if length == 7 {
		y := generic.Parse(params[0]).Int()
		var m = 1
		if v, ok := params[1].(int); ok {
			m = v
		} else if v, ok := params[1].(string); ok {
			m = generic.Parse(v).Int()
		}
		d := generic.Parse(params[2]).Int()
		h := generic.Parse(params[3]).Int()
		i := generic.Parse(params[4]).Int()
		s := generic.Parse(params[5]).Int()
		ns := generic.Parse(params[6]).Int()
		return Time(time.Date(y, time.Month(m), d, h, i, s, ns, time.UTC))
	} else if length == 8 {
		y := generic.Parse(params[0]).Int()
		var m = 1
		if v, ok := params[1].(int); ok {
			m = v
		} else if v, ok := params[1].(string); ok {
			m = generic.Parse(v).Int()
		}
		d := generic.Parse(params[2]).Int()
		h := generic.Parse(params[3]).Int()
		i := generic.Parse(params[4]).Int()
		s := generic.Parse(params[5]).Int()
		ns := generic.Parse(params[6]).Int()
		var location *time.Location
		if v, ok := params[7].(time.Location); ok {
			location = &v
		} else if v, ok := params[7].(*time.Location); ok {
			location = v
		} else {
			location = time.UTC
		}
		return Time(time.Date(y, time.Month(m), d, h, i, s, ns, location))
	} else {
		return Time(time.Now())
	}
}

// GoString implements fmt.GoStringer and formats t to be printed in Go source
func (d Time) GoString() string {
	return d.Time().GoString()
}

// Add returns the time t+d.
func (t Time) Add(d time.Duration) Time {
	return Time(t.Time().Add(d))
}

// AddDate returns the time corresponding to adding the
// given number of years, months, and days to t.
// For example, AddDate(-1, 2, 3) applied to January 1, 2011
// returns March 4, 2010.
//
// AddDate normalizes its result in the same way that Date does,
// so, for example, adding one month to October 31 yields
// December 1, the normalized form for November 31.
func (t Time) AddDate(years int, months int, days int) Time {
	return Time(t.Time().AddDate(years, months, days))
}

// UTC returns t with the location set to UTC.
func (t Time) UTC() Time {
	return Time(t.Time().UTC())
}

// Date returns the year, month, and day in which t occurs.
func (t Time) Date() (year int, month time.Month, day int) {
	return t.Time().Date()
}

// Local returns t with the location set to local time.
func (t Time) Local() Time {
	return Time(t.Time().Local())
}

// Location returns the time zone information associated with t.
func (t Time) Location() *time.Location {
	return t.Time().Location()
}

// Year returns the year in which t occurs.
func (t Time) Year() int {
	return t.Time().Year()
}

// Month returns the month of the year specified by t.
func (t Time) Month() time.Month {
	return t.Time().Month()
}

// Day returns the day of the month specified by t.
func (t Time) Day() int {
	return t.Time().Day()
}

// Hour returns the hour within the day specified by t, in the range [0, 23].
func (t Time) Hour() int {
	return t.Time().Hour()
}

// Minute returns the minute offset within the hour specified by t, in the range [0, 59].
func (t Time) Minute() int {
	return t.Time().Minute()
}

// Second returns the second offset within the minute specified by t, in the range [0, 59].
func (t Time) Second() int {
	return t.Time().Second()
}

// Nanosecond returns the nanosecond offset within the second specified by t,
// in the range [0, 999999999].
func (t Time) Nanosecond() int {
	return t.Time().Nanosecond()
}

// YearDay returns the day of the year specified by t, in the range [1,365] for non-leap years,
// and [1,366] in leap years.
func (t Time) YearDay() int {
	return t.Time().YearDay()
}

// ISOWeek returns the ISO 8601 year and week number in which t occurs.
// Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to
// week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1
// of year n+1.
func (t Time) ISOWeek(d time.Duration) (year, week int) {
	return t.Time().ISOWeek()
}

// Weekday returns the day of the week specified by t.
func (t Time) Weekday() time.Weekday {
	return t.Time().Weekday()
}

// IsDST reports whether the time in the configured location is in Daylight Savings Time.
func (t Time) IsDST() bool {
	return t.Time().IsDST()
}

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (t Time) IsZero() bool {
	return t.Time().IsZero()
}

// After reports whether the time instant t is after u.
func (t Time) After(u interface{}) bool {
	return t.Time().After(Date(u).Time())
}

// Before reports whether the time instant t is before u.
func (t Time) Before(u interface{}) bool {
	return t.Time().Before(Date(u).Time())
}

// Equal reports whether t and u represent the same time instant.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 and 4:00 UTC are Equal.
// See the documentation on the Time type for the pitfalls of using == with
// Time values; most code should use Equal instead.
func (t Time) Equal(u interface{}) bool {
	return t.Time().Equal(Date(u).Time())
}

// UnixMicro returns t as a Unix time, the number of microseconds elapsed since
// January 1, 1970 UTC. The result is undefined if the Unix time in
// microseconds cannot be represented by an int64 (a date before year -290307 or
// after year 294246). The result does not depend on the location associated
// with t.
func (t Time) UnixMicro() int64 {
	return t.Time().UnixMicro()
}

// UnixMilli returns t as a Unix time, the number of milliseconds elapsed since
// January 1, 1970 UTC. The result is undefined if the Unix time in
// milliseconds cannot be represented by an int64 (a date more than 292 million
// years before or after 1970). The result does not depend on the
// location associated with t.
func (t Time) UnixMilli() int64 {
	return t.Time().UnixMilli()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (t *Time) UnmarshalBinary(b []byte) error {
	ti := t.Time()
	var err = ti.UnmarshalBinary(b)
	d := Time(ti)
	t = &d
	return err
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Round operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Round(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t Time) Round(d time.Duration) Time {
	return Time(t.Time().Round(d))
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t Time) Truncate(d time.Duration) Time {
	return Time(t.Time().Truncate(d))
}

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned.
// To compute t-d for a duration d, use t.Add(-d).
func (t Time) Sub(in interface{}) time.Duration {
	return t.Time().Sub(Date(in).Time())
}

// Clock returns the hour, minute, and second within the day specified by t.
func (t Time) Clock() (hour, min, sec int) {
	return t.Time().Clock()
}

// Zone computes the time zone in effect at time t, returning the abbreviated
// name of the zone (such as "CET") and its offset in seconds east of UTC.
func (t Time) Zone() (name string, offset int) {
	return t.Time().Zone()
}

// In returns a copy of t representing the same time instant, but
// with the copy's location information set to loc for display
// purposes.
//
// In panics if loc is nil.
func (t Time) In(loc *time.Location) Time {
	return Time(t.Time().In(loc))
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
func (t *Time) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}

// TODO(rsc): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary, MarshalText,
// UnmarshalBinary, UnmarshalText.

// GobEncode implements the gob.GobEncoder interface.
func (t *Time) GobEncode() ([]byte, error) {
	return t.Time().GobEncode()
}

// GobDecode implements the gob.GobDecoder interface.
func (t *Time) GobDecode(b []byte) error {
	ti := t.Time()
	var err = ti.GobDecode(b)
	d := Time(ti)
	t = &d
	return err
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (t *Time) UnmarshalText(b []byte) error {
	ti := t.Time()
	var err = ti.UnmarshalText(b)
	d := Time(ti)
	t = &d
	return err
}

// Midnight return midnight of given date
//  @return *Date
func (d *Time) Midnight() Time {
	var t = d.Time()
	var time = Time(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()))
	d = &time
	return *d
}

// Calculate calculates relative date to given date
//  @receiver d
//  @param expr
//  @return *Date
//  @return error
func (d *Time) Calculate(expr string) (Time, error) {
	expr = strings.ToLower(expr)
	fields := strings.Fields(expr)
	if len(fields) == 0 {
		return *d, fmt.Errorf("unable to parse date expression:%s", expr)
	}
	if strings.Contains(expr, "midnight") {
		d.Midnight()
	}
	var base = d.Time()
	if fields[0] == "tomorrow" {
		var t = Time(base.AddDate(0, 0, 1))
		d = &t
		return *d, nil
	} else if fields[0] == "yesterday" {
		var t = Time(base.AddDate(0, 0, -1))
		d = &t
		return *d, nil
	} else if fields[0] == "today" {
		return *d, nil
	}
	if len(fields) < 2 {
		return *d, fmt.Errorf("unable to parse date expression:%s", expr)
	}
	var i int
	if fields[0] == "message" {
		i = 1
	} else if fields[0] == "last" {
		i = 2
	} else {
		var err error
		i, err = strconv.Atoi(fields[0])
		if err != nil {
			return *d, fmt.Errorf("unable to parse date expression:%s", expr)
		}
		if len(fields) > 2 {
			if fields[2] == "after" {
				if i < 0 {
					i = i * -1
				}
			}
			if fields[2] == "before" {
				if i > 0 {
					i = i * -1
				}
			}
		}
	}
	var ti time.Time
	if strings.HasPrefix(fields[1], "year") {
		ti = base.AddDate(i, 0, 0)
		if strings.Contains(expr, "start") {
			ti = time.Date(ti.Year(), 1, 1, 0, 0, 0, 0, ti.Location())
		}
	} else if strings.HasPrefix(fields[1], "month") {
		ti = ti.AddDate(0, i, 0)
		if strings.Contains(expr, "start") {
			ti = time.Date(ti.Year(), ti.Month(), 0, 0, 0, 0, 0, ti.Location())
		}
	} else if strings.HasPrefix(fields[1], "day") {
		ti = ti.AddDate(0, 0, i)
		if strings.Contains(expr, "start") {
			d.Midnight()
		}
	} else if strings.HasPrefix(fields[1], "week") {
		ti = ti.AddDate(0, 0, i*7)

		if strings.Contains(expr, "start") {
			// Roll back to Monday:
			if wd := ti.Weekday(); wd == time.Sunday {
				ti = ti.AddDate(0, 0, -6)
			} else {
				ti = ti.AddDate(0, 0, -int(wd)+1)
			}
			d.Midnight()
		}

	} else if strings.HasPrefix(fields[1], "hour") {
		ti = ti.Add(time.Duration(i) * time.Hour)
		if strings.Contains(expr, "start") {
			ti = time.Date(ti.Year(), ti.Month(), ti.Day(), ti.Hour(), 0, 0, 0, ti.Location())
		}
	} else if strings.HasPrefix(fields[1], "minute") {
		ti = ti.Add(time.Duration(i) * time.Minute)
		if strings.Contains(expr, "start") {
			ti = time.Date(ti.Year(), ti.Month(), ti.Day(), ti.Hour(), ti.Minute(), 0, 0, ti.Location())
		}
	} else if strings.HasPrefix(fields[1], "second") {
		ti = ti.Add(time.Duration(i) * time.Second)
	} else {
		return *d, nil
	}
	date := Time(ti)
	d = &date
	return *d, nil

}

// DiffUnix add int64 to given date then return timestamp
//  @receiver d
//  @param t
//  @return time.Duration
func (d *Time) DiffUnix(t int64) time.Duration {
	return time.Duration(d.Unix()-t) * time.Second
}

// DiffDate add date to given date return timestamp\
//  @receiver d
//  @param t
//  @return time.Duration
func (d *Time) DiffDate(t Time) time.Duration {
	return time.Duration(d.Unix()-t.Unix()) * time.Second
}

// DiffExpr add expr to date return timestamp
//  @receiver d
//  @param expr
//  @return time.Duration
//  @return error
func (d *Time) DiffExpr(expr string) (time.Duration, error) {
	t := d.Time()
	_, err := d.Calculate(expr)
	if err != nil {
		return time.Duration(0), err
	}
	return d.DiffTime(t), nil
}

// DiffTime add given time date return timestamp
//  @receiver d
//  @param t
//  @return time.Duration
func (d *Time) DiffTime(t time.Time) time.Duration {
	return time.Duration(d.Unix()-t.Unix()) * time.Second
}

// FormatS format given date as strftime syntax
//  @receiver d
//  @param expr
//  @return string
func (d *Time) FormatS(expr string) string {
	var t = d.Time()
	return strftime.Format(&t, expr)
}

// Format formats given date to given layout as string
//  @receiver d
//  @param layout
//  @param options... interface{}
//  @accept locale as string example:en-US
//  @accept locale as language.Locale
//  @return string
func (d Time) Format(layout string, options ...interface{}) string {
	var t = d.Time()
	var locale = defaultLocale
	for _, item := range options {
		switch option := item.(type) {
		case language.Tag:
			locale = option.String()
		case string:
			locale = GuessLocale(option).String()
		}
	}
	return monday.Format(t, layout, monday.Locale(locale))
}

// Time return time.Time value
//  @receiver d
//  @return time.Time
func (d Time) Time() time.Time {
	var t = time.Time(d)
	return t
}

// Unix return timestamp of given date
//  @receiver d
//  @return int64
func (d *Time) Unix() int64 {
	return d.Unix()
}

// UnixNano return nano timestamp of given date
//  @receiver d
//  @return int64
func (d *Time) UnixNano() int64 {
	return d.UnixNano()
}

// FromString parse any string to Date
//  @param expr
//  @return *Date
//  @return error
func FromString(expr string) (Time, error) {
	t, err := dateparse.ParseLocal(expr)
	if err != nil {
		return Time{}, err
	}
	return Time(t), nil
}

// FromTime parse time to Date
//  @param t
//  @return *Date
func FromTime(t time.Time) Time {
	var time = Time(t)
	return time
}

// FromUnix parse timestamp to Date
//  @param sec
//  @return *Date
func FromUnix(sec int64) Time {
	return Time(time.Unix(sec, 0))
}

// TryParseTime parse anything into Time
//  @param in
//  @return Time
//  @return error
func TryParseTime(in interface{}) (Time, error) {

	return Time{}, fmt.Errorf("unrecognized date input")
}
