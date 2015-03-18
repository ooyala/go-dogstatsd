// Copyright 2013 Ooyala, Inc.

/*
Package dogstatsd provides a Go DogStatsD client. DogStatsD extends StatsD - adding tags and
histograms. Refer to http://docs.datadoghq.com/guides/dogstatsd/ for information about DogStatsD.

Example Usage:
		// Create the client
		c, err := dogstatsd.New("127.0.0.1:8125")
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		// Prefix every metric with the app name
		c.SetGlobalNamespace("flubber.")
		// Send the EC2 availability zone as a tag with every metric
		c.SetGlobalTags([]string{"us-east-1a"})

		err = c.Gauge("request.duration", 1.2, nil, 1)

		// Post info to datadog event stream
		err = c.Info("cookie alert", "Cookies up for grabs in the kitchen!", nil)

dogstatsd is based on go-statsd-client.
*/
package dogstatsd

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type (
	// AlertType represents the supported alert_types of Datadog events.
	AlertType string
	// PriorityType represents Datadog event priority (e.g. 'normal' or 'low')
	PriorityType string
)

// AlertType enum
const (
	Info    AlertType = "info"
	Success AlertType = "success"
	Warning AlertType = "warning"
	Error   AlertType = "error"
)

// PriorityType enum
const (
	Normal PriorityType = "normal"
	Low    PriorityType = "low"
)

const maxEventBytes = 8192

// EventOpts represents detailed options for Event generation
type EventOpts struct {
	DateHappened   time.Time
	Priority       PriorityType
	Host           string
	AggregationKey string
	SourceTypeName string
	Tags           []string
	AlertType      AlertType
}

// Client represents the statsd Client.
type Client struct {
	conn        io.WriteCloser
	eventSource string // eventSource is the Namespace truncated to the first `.`
	namespace   string // Namespace to prepend to all statsd calls
	tags        string // Global tags to be added to every statsd call

	hasNS   bool
	hasTags bool
}

// New returns a pointer to a new Client and an error.
// addr must have the format "hostname:port".
func New(addr string) (*Client, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
	}, nil
}

// send handles sampling and sends the message over UDP. It also adds global namespace prefixes and tags.
func (c *Client) send(name string, value string, tags []string, rate float64) error {
	if rate < 1 {
		// rand.Float64 returns between 0.0 and 1.0. When rate < 1, randomly drop the stat.
		if rand.Float64() > rate {
			return nil
		}
		value += "|@" + strconv.FormatFloat(rate, 'f', 6, 64)
	}

	if c.hasNS {
		name = c.namespace + name
	}

	if c.hasTags {
		tags = append(tags, c.tags)
	}
	if len(tags) > 0 {
		value += "|#" + strings.Join(tags, ",")
	}

	_, err := c.conn.Write([]byte(name + ":" + value))
	return err
}

// Close closes the connection to the DogStatsD agent
func (c *Client) Close() error { return c.conn.Close() }

func (c *Client) newDefaultEventOpts(alertType AlertType, tags []string) *EventOpts {
	return &EventOpts{
		AlertType:      alertType,
		Tags:           tags,
		SourceTypeName: c.eventSource,
	}
}

// Info forges an Info event
func (c *Client) Info(title string, text string, tags []string) error {
	return c.Event(title, text, c.newDefaultEventOpts(Info, tags))
}

// Success forges a Success event
func (c *Client) Success(title string, text string, tags []string) error {
	return c.Event(title, text, c.newDefaultEventOpts(Success, tags))
}

// Warning forges a Warning event
func (c *Client) Warning(title string, text string, tags []string) error {
	return c.Event(title, text, c.newDefaultEventOpts(Warning, tags))
}

// Error forges an Error event
func (c *Client) Error(title string, text string, tags []string) error {
	return c.Event(title, text, c.newDefaultEventOpts(Error, tags))
}

// Event posts to the Datadog event stream.
// Four event types are supported: info, success, warning, error.
// If client Namespace is set it is used as the Event source.
func (c *Client) Event(title string, text string, eo *EventOpts) error {
	// Can't use `len()` because we accept utf8
	titleLen, textLen := utf8.RuneCountInString(title), utf8.RuneCountInString(text)

	eventStr := "_e{" + strconv.FormatInt(int64(titleLen), 10) + "," + strconv.FormatInt(int64(textLen), 10) + "}:" + title + "|" + text + "|t:" + string(eo.AlertType)

	if eo.SourceTypeName != "" {
		eventStr += "|s:" + eo.SourceTypeName
	}
	if !eo.DateHappened.IsZero() {
		eventStr += "|d:" + strconv.FormatInt(eo.DateHappened.Unix(), 10)
	}
	if eo.Priority != "" {
		eventStr += "|p:" + string(eo.Priority)
	}
	if eo.Host != "" {
		eventStr += "|h:" + eo.Host
	}
	if eo.AggregationKey != "" {
		eventStr += "|k:" + eo.AggregationKey
	}
	tags := eo.Tags
	if c.hasTags {
		tags = append(eo.Tags, c.tags)
	}
	if len(tags) > 0 {
		eventStr += "|#" + strings.Join(tags, ",")
	}
	if len(eventStr) > maxEventBytes {
		return fmt.Errorf("Event %q payload is too big (more that 8KB), event discarded", title)
	}
	_, err := c.conn.Write([]byte(eventStr))
	return err
}

// Gauge measures the value of a metric at a particular time
func (c *Client) Gauge(name string, value float64, tags []string, rate float64) error {
	stat := strconv.FormatFloat(value, 'f', 6, 64) + "|g"
	return c.send(name, stat, tags, rate)
}

// Count tracks how many times something happened per second
func (c *Client) Count(name string, value int64, tags []string, rate float64) error {
	stat := strconv.FormatInt(value, 10) + "|c"
	return c.send(name, stat, tags, rate)
}

// Histogram tracks the statistical distribution of a set of values
func (c *Client) Histogram(name string, value float64, tags []string, rate float64) error {
	stat := strconv.FormatFloat(value, 'f', 6, 64) + "|h"
	return c.send(name, stat, tags, rate)
}

// Set counts the number of unique elements in a group
func (c *Client) Set(name string, value string, tags []string, rate float64) error {
	stat := value + "|s"
	return c.send(name, stat, tags, rate)
}

// SetGlobalTags sets the global tags.
func (c *Client) SetGlobalTags(tags []string) {
	c.hasTags = len(tags) != 0
	c.tags = strings.Join(tags, ",")
}

// GetGlobalTags returns the current global tags
func (c *Client) GetGlobalTags() []string {
	return strings.Split(c.tags, ",")
}

// SetGlobalNamespace sets the global namespace and infers the eventSource name
func (c *Client) SetGlobalNamespace(namespace string) {
	c.hasNS = namespace != ""
	c.namespace = namespace
	if period := strings.IndexByte(namespace, '.'); period > -1 {
		c.eventSource = namespace[:period]
	}
}

// GetGlobalNamespace returns the current global namespace
func (c *Client) GetGlobalNamespace() string {
	return c.namespace
}
