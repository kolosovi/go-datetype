package datetype

import (
	"database/sql/driver"
	_ "embed"
	"testing"
	"time"
)

//go:embed tests/tzdata
var tzdata []byte

func TestDate_Time(t *testing.T) {
	nonUTCLoc, err := time.LoadLocationFromTZData("Europe/Moscow", tzdata)
	if err != nil {
		t.Fatalf("could not load location: %v", err)
	}
	tests := []struct {
		name string
		date Date
		loc  *time.Location
		want time.Time
	}{
		{
			name: "null value",
			loc:  time.UTC,
		},
		{
			name: "positive infinity",
			date: Date{Year: 2042, InfinityModifier: Infinity},
			loc:  time.UTC,
		},
		{
			name: "negative infinity",
			date: Date{Year: 2042, InfinityModifier: NegativeInfinity},
			loc:  time.UTC,
		},
		{
			name: "plain value 1",
			date: Date{Year: 2042, Month: time.April, Day: 23, Valid: true},
			loc:  time.UTC,
			want: time.Date(2042, time.April, 23, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "plain value 2",
			date: Date{Year: 2042, Month: time.April, Day: 23, Valid: true},
			loc:  nonUTCLoc,
			want: time.Date(2042, time.April, 23, 0, 0, 0, 0, nonUTCLoc),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.date.Time(tt.loc)
			if got != tt.want {
				t.Errorf("Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Scan(t *testing.T) {
	tests := []struct {
		name    string
		src     interface{}
		wantErr string
		want    Date
	}{
		{
			name: "null",
			want: NewNullDate(),
		},
		{
			name: "infinity",
			src:  "infinity",
			want: NewInfinityDate(),
		},
		{
			name: "-infinity",
			src:  "-infinity",
			want: NewNegativeInfinityDate(),
		},
		{
			name: "plain",
			src:  time.Date(2042, time.April, 23, 1, 2, 3, 4, time.UTC),
			want: NewDate(2042, time.April, 23),
		},
		{
			name:    "unknown type",
			src:     int64(42),
			wantErr: "value 42 has unexpected type int64",
		},
		{
			name:    "unknown string value",
			src:     "sike",
			wantErr: "unexpected string value sike",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var date Date
			err := date.Scan(tt.src)
			if tt.wantErr == "" && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantErr != "" && err == nil {
				t.Errorf("expected an error")
			}
			if tt.wantErr != "" && err.Error() != tt.wantErr {
				t.Errorf("expected error with string \"%v\", got \"%v\"", tt.wantErr, err.Error())
			}
			if err == nil && date != tt.want {
				t.Errorf("scan result = %v, want %v", date, tt.want)
			}
		})
	}
}

func TestDate_Value(t *testing.T) {
	tests := []struct {
		name string
		date Date
		want driver.Value
	}{
		{
			name: "null",
			date: NewNullDate(),
			want: nil,
		},
		{
			name: "infinity",
			date: NewInfinityDate(),
			want: driver.Value("infinity"),
		},
		{
			name: "-infinity",
			date: NewNegativeInfinityDate(),
			want: driver.Value("-infinity"),
		},
		{
			name: "plain",
			date: NewDate(2042, time.April, 23),
			want: driver.Value(time.Date(2042, time.April, 23, 0, 0, 0, 0, time.UTC)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.date.Value()
			if err != nil {
				t.Errorf("got error: %v", err.Error())
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
