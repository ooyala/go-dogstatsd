package dogstatsd

import (
	"fmt"
	"github.com/jabong/blaze/conf"
	"strings"
)

/* Wrapper for passing data to datadog */

var c *Client

func Plot(ploterror int, String string) string {
	if ploterror == 1 {
		s := []string{String, "number_of_errors"}
		return strings.Join(s, ".")
	} else {
		s := []string{String, "number_of_hits"}
		return strings.Join(s, ".")
	}
}

//Initializes StatsD Client and returns a pointer to object
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

//Increments the metric by a defined value
func (c *Client) Increment(name string, value int64, tags []string, rate float64) {
	if c != nil {
		err := c.Incre(name, value, tags, rate)
		if err != nil {
			fmt.Println(err)
		}
	}

}

//Decrements the metric by a defined value
func (c *Client) Decrement(name string, value int64, tags []string, rate float64) {
	if c != nil {
		err := c.Decre(name, value, tags, rate)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Event posts to the Datadog event stream.
// Four event types are supported: info, success, warning, error.
// If client Namespace is set it is used as the Event source.
func (c *Client) Event(title string, text string, eo *EventOpts) {
	if c != nil {
		err := c.Eve(title, text, NewDefaultEventOpts(eo.AlertType, c.Tags, c.Namespace))
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Info posts string with the title to the Datadog event stream
func (c *Client) Information(title string, text string, tags []string) {
	if c != nil {
		err := c.Info(title, text, tags)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Success posts string with the title to the Datadog event stream
func (c *Client) Success(title string, text string, tags []string) {
	if c != nil {
		err := c.Succ(title, text, tags)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Warning posts string with the title to the Datadog event stream
func (c *Client) Warning(title string, text string, tags []string) {
	if c != nil {
		err := c.Warn(title, text, tags)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Error posts string with the title to the Datadog event stream
func (c *Client) Error(title string, text string, tags []string) {
	if c != nil {
		err := c.Err(title, text, tags)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Set counts the number of unique elements in a group
func (c *Client) Set(name string, value string, tags []string, rate float64) {
	if c != nil {
		err := c.Sets(name, value, tags, rate)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Gauge measure the value of a metric at a particular time
func (c *Client) Gauge(name string, value float64, tags []string, rate float64) {
	if c != nil {
		err := c.Gauges(name, value, tags, rate)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}

	}
}

// Histogram tracks the statistical distribution of a set of values
func (c *Client) Histogram(name string, value float64, tags []string, rate float64) {
	if c != nil {
		err := c.Hist(name, value, tags, rate)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Count tracks how many times something happened per second
func (c *Client) Count(name string, value int64, tags []string, rate float64) {
	if c != nil {
		err := c.Counts(name, value, tags, rate)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}
