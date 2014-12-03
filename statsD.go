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
)

// Appends the required suffix to the metric after proper formatting, based on the ploterror parameter
func (c *Client) PlotTitle(ploterror int64, String string) string {
	String = FormatTitle(String)
	if ploterror == 1 {
		s := []string{String, "number_of_errors"}
		return strings.Join(s, ".")
	} else {
		s := []string{String, "number_of_hits"}
		return strings.Join(s, ".")
	}
}

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
func GetFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	if runtime.FuncForPC(pc) != nil {
		function := runtime.FuncForPC(pc).Name()
		return (strings.Split(function, "."))[2]
	}
	return ""

}

// Increments the metric by a defined value, the flag value tells if errors or hits are to be incremented
func (c *Client) Incr(flag int64, value int64, tags []string, rate float64) {
	if c != nil {
		name := c.PlotTitle(flag, GetFunctionName())
		err := c.Increment(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}

}

// Decrements the metric by a defined value, the flag value tells if errors or hits are to be decremented
func (c *Client) Decr(flag int64, value int64, tags []string, rate float64) {
	if c != nil {
		// Finds the name of the calling function which is used as the graph title on Datadog Dashboard
		name := c.PlotTitle(flag, GetFunctionName())
		err := c.Decrement(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}

// Event posts to the Datadog event stream.
// Four event types are supported: info, success, warning, error.
// If client Namespace is set it is used as the Event source.
func (c *Client) Evnt(flag int64, text string, eo *EventOpts) {
	if c != nil {
		title := c.PlotTitle(flag, GetFunctionName())
		err := c.Event(title, text, NewDefaultEventOpts(eo.AlertType, c.Tags, c.Namespace))
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}

// Info posts string with the title to the Datadog event stream
func (c *Client) Inform(flag int64, text string, tags []string) {
	if c != nil {
		title := c.PlotTitle(flag, GetFunctionName())
		err := c.Info(title, text, tags)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}

// Success posts string with the title to the Datadog event stream
func (c *Client) Succss(flag int64, text string, tags []string) {
	if c != nil {
		title := c.PlotTitle(flag, GetFunctionName())
		err := c.Success(title, text, tags)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}

// Warning posts string with the title to the Datadog event stream
func (c *Client) Caution(flag int64, text string, tags []string) {
	if c != nil {
		title := c.PlotTitle(flag, GetFunctionName())
		err := c.Warning(title, text, tags)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}

// Error posts string with the title to the Datadog event stream
func (c *Client) Err(flag int64, text string, tags []string) {
	if c != nil {
		title := c.PlotTitle(flag, GetFunctionName())
		err := c.Error(title, text, tags)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}

// Set counts the number of unique elements in a group
func (c *Client) Sets(flag int64, value string, tags []string, rate float64) {
	if c != nil {
		name := c.PlotTitle(flag, GetFunctionName())
		err := c.Set(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}

// Gauge measure the value of a metric at a particular time
func (c *Client) Gaug(flag int64, value float64, tags []string, rate float64) {
	if c != nil {
		name := c.PlotTitle(flag, GetFunctionName())
		err := c.Gauge(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}

	}
}

// Histogram tracks the statistical distribution of a set of values
func (c *Client) Hist(flag int64, value float64, tags []string, rate float64) {
	if c != nil {
		name := c.PlotTitle(flag, GetFunctionName())
		err := c.Histogram(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}

// Count tracks how many times something happened per second
func (c *Client) Compute(flag int64, value int64, tags []string, rate float64) {
	if c != nil {
		name := c.PlotTitle(flag, GetFunctionName())
		err := c.Count(name, value, tags, rate)
		if err != nil {
			c.slog.Errf("Error: %s", err)
		}
	}
}
