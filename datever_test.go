package datever

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		a := Version{
			Year:     2023,
			Month:    time.June,
			Sequence: 4,
		}
		require.Equal(t, "2023.6.4", a.String())
	})
	t.Run("Increment", func(t *testing.T) {
		a := Version{
			Year:     2023,
			Month:    time.May,
			Sequence: 3,
		}
		b := a.Increment(time.Date(2023, time.June, 3, 0, 30, 0, 0, time.UTC))
		require.Equal(t, "2023.6.0", b.String())
	})
}
