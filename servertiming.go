// Simple library for adding Server-Timing headers
// (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Server-Timing) to
// your application.
//
// Usage:
//
//     Create a new instance:
//
//         t := servertiming.New()
//
//     // Add a few metrics, either by manually specifying the duration:
//
//         t.Add("cache", "Cache Read", 23200*time.Microsecond)
//
//     // Or by using the start-stop API:
//
//         ti.Start("db", "Database query")
//         // query db
//         ti.Stop("db")
//
//     // Then send it with your response:
//
//         w.Header().Set("Server-Timing", ti.String())
//
package servertiming

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Timing struct {
	itemLock sync.Mutex
	items    []*item
}

type item struct {
	Name        string
	Description string
	Duration    time.Duration
	Started     time.Time
}

func New() *Timing {
	return &Timing{}
}

func (t *Timing) String() string {
	t.itemLock.Lock()
	defer t.itemLock.Unlock()

	parts := make([]string, 0)
	for _, item := range t.items {
		subParts := []string{item.Name}
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

func (t *Timing) Add(name, description string, duration time.Duration) {
	i := t.add(name, description)
	i.Duration = duration
}

func (t *Timing) Start(name, description string) {
	i := t.add(name, description)
	i.Started = time.Now()
}

func (t *Timing) Stop(name string) {
	t.itemLock.Lock()
	defer t.itemLock.Unlock()

	for _, item := range t.items {
		if item.Name == name {
			item.Duration = time.Since(item.Started)
		}
	}
}

func (t *Timing) AddFlag(name, description string) {
	t.add(name, description)
}
