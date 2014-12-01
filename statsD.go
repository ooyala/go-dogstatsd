package statsD

import (
	"fmt"
	"github.com/jabong/blaze/conf"
	"github.com/jabong/go-dogstatsd/dogstatsd"
)

/* Wrapper for passing data to datadog */
type Client struct {
	Conn *dogstatsd.Client
}

var c *Client

const (
	ORDER_ADDRESS_NUMBER_OF_HITS            string = "order_address_number_of_hits"
	ORDER_ADDRESS_NUMBER_OF_ERRORS          string = "order_address_number_of_errors"
	CUSTOMER_SEGMENT_NUMBER_OF_HITS         string = "customer_segment_number_of_hits"
	CUSTOMER_SEGMENT_NUMBER_OF_ERRORS       string = "customer_segment_number_of_errors"
	COMPLETE_INFO_NUMBER_OF_HITS            string = "complete_info_number_of_hits"
	COMPLETE_INFO_NUMBER_OF_ERRORS          string = "complete_info_number_of_errors"
	ORDER_LIST_NUMBER_OF_HITS               string = "order_list_number_of_hits"
	ORDER_LIST_NUMBER_OF_ERRORS             string = "order_list_number_of_errors"
	ITEM_CUSTOM_TEXT_NUMBER_OF_HITS         string = "item_custom_text_number_of_hits"
	ITEM_CUSTOM_TEXT_NUMBER_OF_ERRORS       string = "item_custom_text_number_of_errors"
	SHIPPING_PARTNER_AGENT_NUMBER_OF_HITS   string = "shipping_partner_agent_number_of_hits"
	SHIPPING_PARTNER_AGENT_NUMBER_OF_ERRORS string = "shipping_partner_agent_number_of_errors"
)

//Initializes StatsD Client and returns a pointer to object
func Adapter() *Client {
	var err error
	if c == nil {
		c = &Client{}
		// Connecting to the Datadog Agent
		addr := conf.Get("db.ddagent_conn", nil).(string)
		c.Conn, err = dogstatsd.New(addr)
		if err != nil {
			fmt.Println(err)
		} else {
			c.Conn.Tags = nil
			// Prefix every metric with the app name
			c.Conn.Namespace = conf.Get("db.ddagent_Namespace", nil).(string)
		}
	}
	return c
}

//Increments the metric by a defined value
func (c *Client) Increment(name string, value int64, tags []string, rate float64) {
	if c != nil {
		err := c.Conn.Increment(name, value, tags, rate)
		if err != nil {
			fmt.Println(err)
		}
	}

}

//Decrements the metric by a defined value
func (c *Client) Decrement(name string, value int64, tags []string, rate float64) {
	if c != nil {
		err := c.Conn.Decrement(name, value, tags, rate)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Event posts to the Datadog event stream.
// Four event types are supported: info, success, warning, error.
// If client Namespace is set it is used as the Event source.
func (c *Client) Event(title string, text string, eo *dogstatsd.EventOpts) {
	if c != nil {
		err := c.Conn.Event(title, text, dogstatsd.NewDefaultEventOpts(eo.AlertType, c.Conn.Tags, c.Conn.Namespace))
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Info posts string with the title to the Datadog event stream
func (c *Client) Info(title string, text string, tags []string) {
	if c != nil {
		err := c.Conn.Info(title, text, tags)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Success posts string with the title to the Datadog event stream
func (c *Client) Success(title string, text string, tags []string) {
	if c != nil {
		err := c.Conn.Success(title, text, tags)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Warning posts string with the title to the Datadog event stream
func (c *Client) Warning(title string, text string, tags []string) {
	if c.Conn != nil {
		err := c.Conn.Warning(title, text, tags)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Error posts string with the title to the Datadog event stream
func (c *Client) Error(title string, text string, tags []string) {
	if c != nil {
		err := c.Conn.Error(title, text, tags)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Set counts the number of unique elements in a group
func (c *Client) Set(name string, value string, tags []string, rate float64) {
	if c != nil {
		err := c.Conn.Set(name, value, tags, rate)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Gauge measure the value of a metric at a particular time
func (c *Client) Gauge(name string, value float64, tags []string, rate float64) {
	if c != nil {
		err := c.Conn.Gauge(name, value, tags, rate)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}

	}
}

// Histogram tracks the statistical distribution of a set of values
func (c *Client) Histogram(name string, value float64, tags []string, rate float64) {
	if c != nil {
		err := c.Conn.Histogram(name, value, tags, rate)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}

// Count tracks how many times something happened per second
func (c *Client) Count(name string, value int64, tags []string, rate float64) {
	if c != nil {
		err := c.Conn.Count(name, value, tags, rate)
		if err != nil {
			fmt.Println("Failed to connect to Datadog Agent.")
			fmt.Println(err)
		}
	}
}
