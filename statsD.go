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
	HITS       int64   = 0
	ERROR      int64   = 1
	STEPSIZE   int64   = 1
	SAMPLERATE float64 = 1.0
)

// Appends the required suffix to the metric after proper formatting, based on the ploterror parameter
func (c *Client) Plot(ploterror int64, String string) string {
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
	fmt.Println(String)
	if ploterror == 1 {
		s := []string{String, "number_of_errors"}
		return strings.Join(s, ".")
	} else {
		s := []string{String, "number_of_hits"}
		return strings.Join(s, ".")
	}
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
			c.Tags = nil // conf.Get("db.ddagent_tags", nil).([]string)
			// Prefix every metric with the app name
			c.Namespace = conf.Get("db.ddagent_Namespace", nil).(string)
		}
	}
	return c
}

// Increments the metric by a defined value, the flag value tells if errors or hits are to be incremented
func (c *Client) Incre(flag int64, value int64, tags []string, rate float64) {
	if c != nil {
		// Finds the name of the calling function which is used as the graph title on Datadog Dashboard
		pc, _, _, _ := runtime.Caller(1)
		function := runtime.FuncForPC(pc).Name()
		fname := strings.Split(function, ".")
		name := c.Plot(flag, fname[2])
		err := c.Increment(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Failed to connect to Datadog Agent.", err)
		}
	}

}

// Decrements the metric by a defined value, the flag value tells if errors or hits are to be decremented
func (c *Client) Decre(flag int64, value int64, tags []string, rate float64) {
	if c != nil {
		// Finds the name of the calling function which is used as the graph title on Datadog Dashboard
		pc, _, _, _ := runtime.Caller(1)
		function := runtime.FuncForPC(pc).Name()
		fname := strings.Split(function, ".")
		name := c.Plot(flag, fname[2])
		err := c.Decrement(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Failed to connect to Datadog Agent %s", err)
		}
	}
}

// Event posts to the Datadog event stream.
// Four event types are supported: info, success, warning, error.
// If client Namespace is set it is used as the Event source.
func (c *Client) Evnt(title string, text string, eo *EventOpts) {
	if c != nil {
		err := c.Event(title, text, NewDefaultEventOpts(eo.AlertType, c.Tags, c.Namespace))
		if err != nil {
			c.slog.Errf("Failed to connect to Datadog Agent.", err)
		}
	}
}

// Info posts string with the title to the Datadog event stream
func (c *Client) Infor(title string, text string, tags []string) {
	if c != nil {
		err := c.Info(title, text, tags)
		if err != nil {
			c.slog.Errf("Failed to connect to Datadog Agent.", err)
		}
	}
}

// Success posts string with the title to the Datadog event stream
func (c *Client) Succss(title string, text string, tags []string) {
	if c != nil {
		err := c.Success(title, text, tags)
		if err != nil {
			c.slog.Errf("Failed to connect to Datadog Agent.", err)
		}
	}
}

// Warning posts string with the title to the Datadog event stream
func (c *Client) Warn(title string, text string, tags []string) {
	if c != nil {
		err := c.Warning(title, text, tags)
		if err != nil {
			c.slog.Errf("Failed to connect to Datadog Agent.", err)
		}
	}
}

// Error posts string with the title to the Datadog event stream
func (c *Client) Err(title string, text string, tags []string) {
	if c != nil {
		err := c.Error(title, text, tags)
		if err != nil {
			c.slog.Errf("Failed to connect to Datadog Agent.", err)
		}
	}
}

// Set counts the number of unique elements in a group
func (c *Client) Sets(name string, value string, tags []string, rate float64) {
	if c != nil {
		err := c.Set(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Failed to connect to Datadog Agent.", err)
		}
	}
}

// Gauge measure the value of a metric at a particular time
func (c *Client) Gauges(name string, value float64, tags []string, rate float64) {
	if c != nil {
		err := c.Gauge(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Failed to connect to Datadog Agent.", err)
		}

	}
}

// Histogram tracks the statistical distribution of a set of values
func (c *Client) Hist(name string, value float64, tags []string, rate float64) {
	if c != nil {
		err := c.Histogram(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Failed to connect to Datadog Agent.", err)
		}
	}
}

// Count tracks how many times something happened per second
func (c *Client) Counts(name string, value int64, tags []string, rate float64) {
	if c != nil {
		err := c.Count(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Failed to connect to Datadog Agent.", err)
		}
	}
}
