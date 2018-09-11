// Simple library for adding Server-Timing headers
// (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Server-Timing) to
// your application.
//
// Usage
//
// A Timing object can be used to gather metrics and format them into the correct format.
//
// Create a new instance:
//
//     t := servertiming.New()
//
// Optionally enable name prefixing to preserve the order of
// metrics (will adjust names though!)
//
//     t.EnablePrefix()
//
// Add a few metrics, either by manually specifying the duration:
//
//     t.Add("cache", "Cache Read", 23200*time.Microsecond)
//
// Or by using the start-stop API:
//
//     ti.Start("db", "Database query")
//     // query db
//     ti.Stop("db")
//
// Then send it with your response:
//
//     w.Header().Set("Server-Timing", ti.String())
//
// HTTP2 Trailers
//
// Timings can be sent as a trailer when using HTTP2, see the example in net/http: https://godoc.org/net/http#example-ResponseWriter--Trailers
package servertiming

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

// Timing holds timing metrics
type Timing struct {
	itemLock sync.Mutex
	items    []*item
	prefix   bool
}

type item struct {
	Name        string
	Description string
	Duration    time.Duration
	Started     time.Time
}

// Creates a new Timing object
func New() *Timing {
	return &Timing{}
}

// Add a numerical prefix to each metric, to preserve ordering when sorted in dev tools.
func (t *Timing) EnablePrefix() *Timing {
	t.prefix = true
	return t
}

// Formats a valid Server-Timing header value, as defined in https://w3c.github.io/server-timing/#the-server-timing-header-field
func (t *Timing) String() string {
	t.itemLock.Lock()
	defer t.itemLock.Unlock()

	pos := int(math.Ceil(math.Log10(float64(len(t.items)))))
	nameFmt := fmt.Sprintf("%%0%dd_%%s", pos)
	parts := make([]string, 0)
	for idx, item := range t.items {
		subParts := []string{}
		if t.prefix {
			subParts = append(subParts, fmt.Sprintf(nameFmt, idx, item.Name))
		} else {
			subParts = append(subParts, item.Name)
		}
		if item.Description != "" {
			subParts = append(subParts, fmt.Sprintf("desc=%#v", item.Description))
		}
		if item.Duration != 0 {
			subParts = append(subParts, fmt.Sprintf("dur=%.2f", item.Duration.Seconds()*1000))
		}

		parts = append(parts, strings.Join(subParts, ";"))
	}

	return strings.Join(parts, ", ")
}

func (t *Timing) add(name, description string) *item {
	t.itemLock.Lock()
	defer t.itemLock.Unlock()

	i := &item{
		Name:        name,
		Description: description,
	}
	t.items = append(t.items, i)

	return i
}

// Add a new timing, with the duration specified
func (t *Timing) Add(name, description string, duration time.Duration) {
	i := t.add(name, description)
	i.Duration = duration
}

// Start a timer
func (t *Timing) Start(name, description string) {
	i := t.add(name, description)
	i.Started = time.Now()
}

// Stop a timer
func (t *Timing) Stop(name string) {
	t.itemLock.Lock()
	defer t.itemLock.Unlock()

	for _, item := range t.items {
		if item.Name == name {
			item.Duration = time.Since(item.Started)
		}
	}
}

// Add a flag without value
func (t *Timing) AddFlag(name, description string) {
	t.add(name, description)
}
