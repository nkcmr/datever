package datever // import "code.nkcmr.net/datever"

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Version represents a version that is formatted as "YYYY.MM.N".
// Example "2023.4.2"
type Version struct {
	Year     int
	Month    time.Month
	Sequence uint
}

func Parse(s string) (Version, error) {
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("expected 3 parts delimited by a period, got %d", len(parts))
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, fmt.Errorf("invalid year integer: %s", err.Error())
	}
	monthInt, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, fmt.Errorf("invalid month integer: %s", err.Error())
	}
	month := time.Month(monthInt)
	if strings.HasPrefix(month.String(), "%!") {
		return Version{}, fmt.Errorf("invalid month: must be in range [1-12]")
	}
	seq, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return Version{}, fmt.Errorf("invalid sequence integer: %s", err.Error())
	}
	return Version{
		Year:     year,
		Month:    month,
		Sequence: uint(seq),
	}, nil
}

// String implements fmt.Stringer
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Year, int(v.Month), v.Sequence)
}

// Increment will increment the current version according to the rules
// of datever.
func (v Version) Increment(now time.Time) Version {
	nowDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	if v.toDate().After(nowDate) {
		return Version{
			Year:     v.Year,
			Month:    v.Month,
			Sequence: v.Sequence + 1,
		}
	}
	seq := v.Sequence + 1
	if v.Year != now.Year() || v.Month != now.Month() {
		seq = 0
	}
	return Version{
		Year:     now.Year(),
		Month:    now.Month(),
		Sequence: seq,
	}
}

func (v Version) toDate() time.Time {
	return time.Date(v.Year, v.Month, 1, 0, 0, 0, 0, time.UTC)
}

// Compare will compare two datever versions. Returning
func Compare(a, b Version) int {
	adate := a.toDate()
	bdate := b.toDate()
	if adate.After(bdate) {
		return 1
	} else if bdate.After(adate) {
		return -1
	}
	if !adate.Equal(bdate) {
		panic("expected dates to be equal if they were not before or after one another")
	}
	if a.Sequence > b.Sequence {
		return 1
	} else if a.Sequence < b.Sequence {
		return -1
	}
	if a.Sequence != b.Sequence {
		panic("expected sequences to be equal if they were not greater or less-than one another")
	}
	return 0
}
