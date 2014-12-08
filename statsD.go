package dogstatsd

import (
	"fmt"
	"github.com/jabong/blaze/conf"
	"runtime"
	"strings"
	"unicode"
)

/* Wrapper for passing data to datadog */

var c *Client

const (
	HITS        int64   = 0
	ERROR       int64   = 1
	STEP_SIZE   int64   = 1
	SAMPLE_RATE float64 = 1.0
	SKIP        int     = 2
)

// Appends the required suffix to the metric after proper formatting, based on the ploterror parameter
func (c *Client) MetricTitle(ploterror int64, String string) string {
	String = FormatTitle(String)
	s := make([]string, 2)
	s[0] = String

	if ploterror == 1 {
		if GetFunctionName(SKIP) != "Incr" && GetFunctionName(SKIP) != "Decr" {
			s[1] = strings.ToLower(GetFunctionName(SKIP))
		} else {
			s[1] = "number_of_errors"
		}
	} else {
		s[1] = "number_of_hits"
	}
	fmt.Println(strings.Join(s, "."))
	return strings.Join(s, ".")
}

// Formats and constructs the title for the metrics and events to be posted on the Datadog API.
func FormatTitle(String string) string {
	var words []string
	l := 0
	for s := String; s != ""; s = s[l:] {
		l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if l <= 0 {
			l = len(s)
		}
		words = append(words, strings.ToLower(s[:l]))
	}
	String = strings.Join(words, "_")
	return String
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
	return c
}

// Finds the name of the calling function which is used as the graph title on Datadog Dashboard
func GetFunctionName(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	if runtime.FuncForPC(pc) != nil {
		function := runtime.FuncForPC(pc).Name()
		return (strings.Split(function, "."))[2]
	}
	return ""
}

// Increments the metric by a defined value, the flag value tells if errors or hits are to be incremented
func (c *Client) Incr(flag int64, value int64, tags []string, rate float64) {
	if c != nil {
		// Gets the calling function's name
		pc, _, _, _ := runtime.Caller(1)
		function := runtime.FuncForPC(pc).Name()
		// Uses the value to get the metric-name,to be graphed on the Datadog dashboard.
		name := c.MetricTitle(flag, (strings.Split(function, "."))[2])
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
		name := c.MetricTitle(flag, (strings.Split(function, "."))[2])
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
		title := c.MetricTitle(flag, (strings.Split(function, "."))[2])
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
		title := c.MetricTitle(flag, (strings.Split(function, "."))[2])
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
		title := c.MetricTitle(flag, (strings.Split(function, "."))[2])
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
		title := c.MetricTitle(flag, (strings.Split(function, "."))[2])
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
		pc, _, _, _ := runtime.Caller(1)
		function := runtime.FuncForPC(pc).Name()
		// Uses the value as title for the Datadog Event Stream.
		title := c.MetricTitle(flag, (strings.Split(function, "."))[2])
		err := c.Error(title, text, tags)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}

// Sets counts the number of unique elements in a group.
func (c *Client) Sets(flag int64, value string, tags []string, rate float64) {
	if c != nil {
		// Gets the calling function's name
		pc, _, _, _ := runtime.Caller(1)
		function := runtime.FuncForPC(pc).Name()
		// Uses the value as title for the Datadog Event Stream.
		name := c.MetricTitle(flag, (strings.Split(function, "."))[2])
		err := c.Set(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}

// Gaug measure the value of a metric at a particular time.
func (c *Client) Gaug(flag int64, value float64, tags []string, rate float64) {
	if c != nil {
		// Gets the calling function's name
		pc, _, _, _ := runtime.Caller(1)
		function := runtime.FuncForPC(pc).Name()
		// Uses the value as title for the Datadog Event Stream.
		name := c.MetricTitle(flag, (strings.Split(function, "."))[2])
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
		name := c.MetricTitle(flag, (strings.Split(function, "."))[2])
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
		name := c.MetricTitle(flag, (strings.Split(function, "."))[2])
		err := c.Count(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}
