package datetype

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Date represents a date that can optionally be null or be equal to positive
// or negative infinity.
type Date struct {
	Year  int
	Month time.Month
	Day   int

	// Is false if and only if date is null.
	Valid bool
	// Is None if and only if date is equal to positive or negative infinity.
	InfinityModifier InfinityModifier
}

// NewNullDate returns a null date.
func NewNullDate() Date {
	return Date{Year: 1, Month: time.January, Day: 1}
}

// NewDate a date with specified year, month & day.
func NewDate(year int, month time.Month, day int) Date {
	return Date{Year: year, Month: month, Day: day, Valid: true}
}

// NewInfinityDate returns a date equal to infinity.
func NewInfinityDate() Date {
	return Date{
		Year: 1,
		Month: time.January,
		Day: 1,
		Valid: true,
		InfinityModifier: Infinity,
	}
}

// NewNegativeInfinityDate returns a date equal to negative infinity.
func NewNegativeInfinityDate() Date {
	return Date{
		Year: 1,
		Month: time.January,
		Day: 1,
		Valid: true,
		InfinityModifier: NegativeInfinity,
	}
}

// InfinityModifier represents a date's infinity status.
type InfinityModifier int8

const (
	// NegativeInfinity means negative infinity.
	NegativeInfinity InfinityModifier = iota - 1
	// None means that a date is not equal to infinity or negative infinity.
	None
	// Infinity means positive infinity.
	Infinity
)

// Time returns time corresponding to this date in timezone loc.
func (d Date) Time(loc *time.Location) time.Time {
	if !d.Valid || d.InfinityModifier != None {
		return time.Date(1, time.January, 1, 0, 0, 0, 0, loc)
	}
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, loc)
}

// Scan implements database/sql Scanner interface.
func (d *Date) Scan(src interface{}) error {
	d.Year = 1
	d.Month = time.January
	d.Day = 1
	if src == nil {
		return nil
	}
	switch srcValue := src.(type) {
	case time.Time:
		d.Year = srcValue.Year()
		d.Month = srcValue.Month()
		d.Day = srcValue.Day()
		d.Valid = true
		return nil
	case string:
		switch srcValue {
		case infinityStr:
			d.Valid = true
			d.InfinityModifier = Infinity
			return nil
		case negativeInfinityStr:
			d.Valid = true
			d.InfinityModifier = NegativeInfinity
			return nil
		}
		return fmt.Errorf("unexpected string value %v", srcValue)
	default:
		return fmt.Errorf("value %#v has unexpected type %T", srcValue, srcValue)
	}
}

// Value implements database/sql Valuer interface.
func (d Date) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	if d.InfinityModifier == Infinity {
		return infinityStr, nil
	}
	if d.InfinityModifier == NegativeInfinity {
		return negativeInfinityStr, nil
	}
	return d.Time(time.UTC), nil
}

const (
	infinityStr         = "infinity"
	negativeInfinityStr = "-infinity"
)
