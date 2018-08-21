package servertiming_test

import (
	"time"

	"github.com/rubenv/servertiming"
)

func Example() {
	ti := servertiming.New().EnablePrefix()

	ti.AddFlag("missedCache", "Cache missed")

	ti.Add("cache", "Cache Read", 23200*time.Microsecond)

	ti.Start("db", "Database query")
	// query db
	ti.Stop("db")

	// Send in response:
	//w.Header().Set("Server-Timing", ti.String())
}
