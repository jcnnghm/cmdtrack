package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Command represents a command that was executed
type Command struct {
	Command    string `json:"Command"`
	Hostname   string `json:"Hostname"`
	WorkingDir string `json:"WorkingDir"`
	Timestamp  int64  `json:"Timestamp"`
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

// FetchCommands fetches history from the CommandTrack server
func FetchCommands(cmdtrackURL string, verbose bool) (commands []Command, err error) {
	commands = make([]Command, 0, 10000)

	if verbose {
		fmt.Println("Starting request to fetch commands")
	}
	req, err := http.NewRequest("GET", cmdtrackURL+"history", nil)
	if err != nil {
		return
	}

	req.Header.Add("Secret", Config.SharedSecret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if verbose {
		fmt.Println("Request complete, parsing...")
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&commands)
	resp.Body.Close()
	if err != nil {
		return
	}
	if verbose {
		fmt.Println("Parsing complete, decrypting...")
	}

	for i := range commands {
		if err = commands[i].decrypt(); err != nil {
			return
		}
	}

	return
}

func (c *Command) decrypt() error {
	cmd, err := DecryptBase64(c.Command, Config.EncryptionKey)
	c.Command = cmd
	return err
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
		if resp, err := postForm(cmdtrackURL+"command", values); err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		time.Sleep(1000 * time.Millisecond)
		count = count + 1
	}
	return errors.New("Failed to save")
}

func postForm(url string, values url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Secret", Config.SharedSecret)
	resp, err = http.DefaultClient.Do(req)
	return
}
