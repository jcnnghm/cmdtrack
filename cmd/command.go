package cmd

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Command represents a command that was executed
type Command struct {
	Command,
	Hostname,
	WorkingDir string
	Timestamp int64
}

// NewCommand creates a command from a request
func NewCommand(r *http.Request) (*Command, error) {
	if timestamp, err := strconv.ParseInt(r.FormValue("Timestamp"), 10, 64); err != nil {
		return nil, err
	} else {
		return &Command{
			Command:    r.FormValue("Command"),
			Hostname:   r.FormValue("Hostname"),
			WorkingDir: r.FormValue("WorkingDir"),
			Timestamp:  timestamp,
		}, nil
	}
}

// Normalize normalizes all of the fields in a command
func (c *Command) Normalize() {
	c.Command = strings.TrimSpace(c.Command)
	c.Hostname = strings.TrimSpace(c.Hostname)
	c.WorkingDir = strings.TrimSpace(c.WorkingDir)

	if len(c.Hostname) == 0 {
		if host, err := os.Hostname(); err == nil {
			c.Hostname = host
		}
	}

	if c.Timestamp == 0 {
		c.Timestamp = time.Now().Unix()
	}
}

// IsValid normalizes a command and verifies that it's valid
func (c *Command) IsValid() bool {
	c.Normalize()
	return len(c.Command) > 0 && len(c.Hostname) > 0 && len(c.WorkingDir) > 0
}

func (c *Command) toURLValues() url.Values {
	return url.Values{"Command": {EncryptBase64(c.Command, Config.EncryptionKey)}, "Hostname": {c.Hostname}, "WorkingDir": {c.WorkingDir}, "Timestamp": {strconv.FormatInt(c.Timestamp, 10)}}
}

// Send sends the command to the cmdtrack server
func (c *Command) Send(cmdtrackURL string) error {
	values := c.toURLValues()
	retryCount := 10
	count := 0
	for count < retryCount {
		if resp, err := http.PostForm(cmdtrackURL+"command", values); err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		time.Sleep(1000 * time.Millisecond)
		count = count + 1
	}
	return errors.New("Failed to save")
}
