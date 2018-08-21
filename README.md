# servertiming

> Go library for adding [`Server-Timing` headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Server-Timing) to your application.

[![Build Status](https://travis-ci.org/rubenv/servertiming.svg?branch=master)](https://travis-ci.org/rubenv/servertiming) [![GoDoc](https://godoc.org/github.com/rubenv/servertiming?status.png)](https://godoc.org/github.com/rubenv/servertiming)

## Usage

```go
// Create a new instance:

    t := servertiming.New()

// Add a few metrics, either by manually specifying the duration:

    t.Add("cache", "Cache Read", 23200*time.Microsecond)

// Or by using the start-stop API:

    ti.Start("db", "Database query")
    // query db
    ti.Stop("db")

// Then send it with your response:

    w.Header().Set("Server-Timing", ti.String())
```

## License

This library is distributed under the [MIT](LICENSE) license.
