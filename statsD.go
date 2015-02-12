package dogstatsd

import (
    "fmt"
    "github.com/jabong/blaze/conf"
    "runtime"
    "strings"
)

/* Wrapper for passing data to datadog */

var c *Client

const (
    HITS        int64   = 0
    ERROR       int64   = 1
    STEP_SIZE   int64   = 1
    SAMPLE_RATE float64 = 1.0
)

// Appends the required suffix to the metric after proper formatting, based on the ploterror parameter
// ploterror is a boolean parameter which is used to check, weather to append "number_of_error" or
// "number_of_hits" to the Title. If the calling function is not "Incr" or "Decr", append the name of
// the function as suffix to the Title.
func (c *Client) MetricTitle(ploterror int64, title string) string {
    s := make([]string, 2)
    s[0] = title
    pc, _, _, _ := runtime.Caller(1)
    function := runtime.FuncForPC(pc).Name()
    fname := (strings.Split(function, "/"))[len(strings.Split(function, "/"))-1]
    fname = (strings.Split(fname, "."))[len(strings.Split(fname, "."))-1]
    if fname != "Incr" && fname != "Decr" {
        s[1] = strings.ToLower(fname)
    } else {
        if ploterror == ERROR {
            s[1] = "number_of_errors"
        } else if ploterror == HITS {
            s[1] = "number_of_hits"
        }
    }
    return strings.Join(s, ".")
}

// Initializes StatsD Client and returns a pointer to object
func Adapter() *Client {
    var err error
    if c == nil {
        c = &Client{}
        // Connecting to the Datadog Agent
        addr := conf.Get("db.ddagent_conn", nil).(string)

        c, err = New(addr)
        if err != nil {
            fmt.Println(err)
        } else {
            c.Tags = nil
            // Prefix every metric with the app name
            c.Namespace = conf.Get("db.ddagent_Namespace", nil).(string)
        }
    }
    c.AllocatedMemory()
    c.NextGC()
    return c
}

// Extracts the function name from the given function path.
func parse(functionName string) string {
    fname := (strings.Split(functionName, "/"))[len(strings.Split(functionName, "/"))-1]
    return fname
}

// Increments the metric by a defined value, the flag value tells if errors or hits are to be incremented
func (c *Client) Incr(flag int64, value int64, tags []string, rate float64) {
    if c != nil {
        // Gets the calling function's name
        pc, _, _, _ := runtime.Caller(1)
        function := runtime.FuncForPC(pc).Name()
        // Uses the value to get the metric-name,to be graphed on the Datadog dashboard.
        name := c.MetricTitle(flag, parse(function))
        err := c.Increment(name, value, tags, rate)
        if err != nil {
            c.slog.Errf("Error: %s", err)
        }
    }
}

// Decrements the metric by a defined value, the flag value tells if errors or hits are to be decremented
func (c *Client) Decr(flag int64, value int64, tags []string, rate float64) {
    if c != nil {
        // Gets the calling function's name
        pc, _, _, _ := runtime.Caller(1)
        function := runtime.FuncForPC(pc).Name()
        // Uses the value to get the metric-name,to be graphed on the Datadog dashboard.
        name := c.MetricTitle(flag, parse(function))
        err := c.Decrement(name, value, tags, rate)
        if err != nil {
            c.slog.Errf("Error: %s", err)
        }
    }
}

// Evnt posts to the Datadog event stream.
// Four event types are supported: info, success, warning, error.
// If client Namespace is set it is used as the Event source.
func (c *Client) Evnt(flag int64, text string, eo *EventOpts) {
    if c != nil {
        // Gets the calling function's name
        pc, _, _, _ := runtime.Caller(1)
        function := runtime.FuncForPC(pc).Name()
        // Uses the value as title for the Datadog Event Stream.
        title := c.MetricTitle(flag, parse(function))
        err := c.Event(title, text, NewDefaultEventOpts(eo.AlertType, c.Tags, c.Namespace))
        if err != nil {
            c.slog.Errf("Error: %s", err)
        }
    }
}

// Inform posts string with event-title as the name of the calling function, to the Datadog event stream
func (c *Client) Inform(flag int64, text string, tags []string) {
    if c != nil {
        // Gets the calling function's name
        pc, _, _, _ := runtime.Caller(1)
        function := runtime.FuncForPC(pc).Name()
        // Uses the value as title for the Datadog Event Stream.
        title := c.MetricTitle(flag, parse(function))
        err := c.Info(title, text, tags)
        if err != nil {
            c.slog.Errf("Error: %s", err)
        }
    }
}

// Succss posts string with event-title as the name of the calling function, to the Datadog event stream
func (c *Client) Succss(flag int64, text string, tags []string) {
    if c != nil {
        // Gets the calling function's name
        pc, _, _, _ := runtime.Caller(1)
        function := runtime.FuncForPC(pc).Name()
        // Uses the value as title for the Datadog Event Stream.
        title := c.MetricTitle(flag, parse(function))
        err := c.Success(title, text, tags)
        if err != nil {
            c.slog.Errf("Error: %s", err)
        }
    }
}

// Caution posts string with event-title as the name of the calling function, to the Datadog event stream
func (c *Client) Caution(flag int64, text string, tags []string) {
    if c != nil {
        // Gets the calling function's name
        pc, _, _, _ := runtime.Caller(1)
        function := runtime.FuncForPC(pc).Name()
        // Uses the value as title for the Datadog Event Stream.
        title := c.MetricTitle(flag, parse(function))
        err := c.Warning(title, text, tags)
        if err != nil {
            c.slog.Errf("Error: %s", err)
        }
    }
}

// Fault posts string with event-title as the name of the calling function, to the Datadog event stream
func (c *Client) Fault(flag int64, text string, tags []string) {
    if c != nil {
        // Gets the calling function's name
        pc, _, linenumber, _ := runtime.Caller(1)
        function := runtime.FuncForPC(pc).Name()
        text = fmt.Sprintf("%s Line Number : %d", text, linenumber)
        // Uses the value as title for the Datadog Event Stream.
        title := c.MetricTitle(flag, parse(function))
        err := c.Error(title, text, tags)
        if err != nil {
            c.slog.Errf("Error: %s", err)
        }
    }
}

// Sets counts the number of unique elements in a group.
func (c *Client) Sets(flag int64, value string, tags []string, rate float64) {
    if c != nil {
        // Gets the cal string,ling function's name
        pc, _, _, _ := runtime.Caller(1)
        function := runtime.FuncForPC(pc).Name()
        // Uses the value as title for the Datadog Event Stream.
        name := c.MetricTitle(flag, parse(function))
        err := c.Set(name, value, tags, rate)
        if err != nil {
            c.slog.Errf("Error: %s", err)
        }
    }
}

// Gaug measure the value of a metric at a particular time.
func (c *Client) Assess(flag int64, value float64, tags []string, rate float64) {
    if c != nil {
        // Gets the calling function's name
        pc, _, _, _ := runtime.Caller(1)
        function := runtime.FuncForPC(pc).Name()
        // Uses the value as title for the Datadog Event Stream.
        name := c.MetricTitle(flag, parse(function))
        err := c.Gauge(name, value, tags, rate)
        if err != nil {
            c.slog.Errf("Error: %s", err)
        }
    }
}

// Hist tracks the statistical distribution of a set of values.
func (c *Client) Hist(flag int64, value float64, tags []string, rate float64) {
    if c != nil {
        // Gets the calling function's name
        pc, _, _, _ := runtime.Caller(1)
        function := runtime.FuncForPC(pc).Name()
        // Uses the value as title for the Datadog Event Stream.
        name := c.MetricTitle(flag, parse(function))
        err := c.Histogram(name, value, tags, rate)
        if err != nil {
            c.slog.Errf("Error: %s", err)
        }
    }
}

// Compute tracks how many times something happened per second
func (c *Client) Compute(flag int64, value int64, tags []string, rate float64) {
    if c != nil {
        // Gets the calling function's name
        pc, _, _, _ := runtime.Caller(1)
        function := runtime.FuncForPC(pc).Name()
        // Uses the value as title for the Datadog Event Stream.
        name := c.MetricTitle(flag, parse(function))
        err := c.Count(name, value, tags, rate)
        if err != nil {
            c.slog.Errf("Error: %s", err)
        }
    }
}
func checkerror(err error) {
    if err != nil {
        c.slog.Errf("Error: %s", err)
    }
}

// Memory in use metrics
// Plots the number of bytes allocated and still in use.
func (c *Client) AllocatedMemory() {
    if c != nil {
        runtime.ReadMemStats(c.MemoryStats)
        title := "allocated_memory"
        err := c.Gauge(title, float64(c.MemoryStats.Alloc), nil, SAMPLE_RATE)
        checkerror(err)
    }
}

// Plots the number of bytes used now by mcache structures
func (c *Client) MCacheInuse() {
    if c != nil {
        runtime.ReadMemStats(c.MemoryStats)
        title := "mcache_inuse"
        err := c.Gauge(title, float64(c.MemoryStats.MCacheInuse), nil, SAMPLE_RATE)
        checkerror(err)
    }
}

// Plots the number of bytes used now by mspan structures
func (c *Client) MSpanInuse() {
    if c != nil {
        runtime.ReadMemStats(c.MemoryStats)
        title := "mspan_inuse"
        err := c.Gauge(title, float64(c.MemoryStats.MSpanInuse), nil, SAMPLE_RATE)
        checkerror(err)
    }
}

// Plots the number of bytes in non-idle span
func (c *Client) HeapInuse() {
    if c != nil {
        runtime.ReadMemStats(c.MemoryStats)
        title := "heap_inuse"
        err := c.Gauge(title, float64(c.MemoryStats.HeapInuse), nil, SAMPLE_RATE)
        checkerror(err)
    }
}

// Sytem memory allocations

// Plots the number of bytes obtained from system
func (c *Client) HeapSystem() {
    if c != nil {
        runtime.ReadMemStats(c.MemoryStats)
        title := "heap_allocation"
        err := c.Gauge(title, float64(c.MemoryStats.HeapSys), nil, SAMPLE_RATE)
        checkerror(err)
    }
}

func (c *Client) MSpanSystem() {
    if c != nil {
        runtime.ReadMemStats(c.MemoryStats)
        title := "mspan_system"
        err := c.Gauge(title, float64(c.MemoryStats.MSpanSys), nil, SAMPLE_RATE)
        checkerror(err)
    }
}

func (c *Client) MCacheSystem() {
    if c != nil {
        runtime.ReadMemStats(c.MemoryStats)
        title := "mcache_system"
        err := c.Gauge(title, float64(c.MemoryStats.MCacheSys), nil, SAMPLE_RATE)
        checkerror(err)
    }
}

// Garbage collector statistics
// next collection will happen when HeapAlloc ≥ this amount
func (c *Client) NextGC() {
    if c != nil {
        runtime.ReadMemStats(c.MemoryStats)
        title := "next garbage collection"
        err := c.Gauge(title, float64(c.MemoryStats.NextGC), nil, SAMPLE_RATE)
        checkerror(err)
    }
}

// Displays the end time of last collection (nanoseconds since 1970)
func (c *Client) GetLastGC() {
    if c != nil {
        runtime.ReadMemStats(c.MemoryStats)
        title := "Last Garbage Collection"
        err := c.Info(title, string(c.MemoryStats.LastGC), nil)
        checkerror(err)
    }
}
