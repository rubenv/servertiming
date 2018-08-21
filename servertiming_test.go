package servertiming

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cheekybits/is"
)

func TestString(t *testing.T) {
	is := is.New(t)

	ti := New()
	is.NotNil(ti)

	ti.AddFlag("missedCache", "")
	is.Equal(ti.String(), "missedCache")

	ti = New()
	ti.AddFlag("missedCache", "Cache missed")
	is.Equal(ti.String(), "missedCache;desc=\"Cache missed\"")

	ti = New()
	ti.AddFlag("missedCache", "Cache missed: \"3\"")
	is.Equal(ti.String(), `missedCache;desc="Cache missed: \"3\""`)

	ti = New()
	ti.AddFlag("missedCache", "")
	ti.Add("cache", "Cache Read", 23200*time.Microsecond)
	is.Equal(ti.String(), `missedCache, cache;desc="Cache Read";dur=23.20`)

	ti = New()
	ti.Start("db", "Database")
	time.Sleep(10 * time.Millisecond)
	ti.Stop("db")
	s := ti.String()
	parts := strings.Split(s, ";")
	is.Equal(len(parts), 3)
	is.Equal(parts[0], "db")
	is.Equal(parts[1], `desc="Database"`)

	tparts := strings.Split(parts[2], "=")
	is.Equal(len(tparts), 2)
	is.Equal(tparts[0], "dur")

	d, err := strconv.ParseFloat(tparts[1], 64)
	is.NoErr(err)
	is.True(d > 10)
}
