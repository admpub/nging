package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/webx-top/echo"
)

// LogEvent struct that is also the Based struct.
type LogEvent struct {
	Timestamp time.Time  `json:"timestamp"`
	Message   string     `json:"message"`
	Tags      []string   `json:"tags,omitempty"`
	Extra     echo.Store `json:"-"`
}

// Formate-Type
var (
	retime = regexp.MustCompile(`%{\+([^}]+)}`)
	revar  = regexp.MustCompile(`%{([\w@]+)}`)
)

const timeFormat = `2006-01-02T15:04:05.999999999Z`

// AddTag for LogEvent Tags
func (le *LogEvent) AddTag(tags ...string) {
	for _, tag := range tags {
		ftag := le.Format(tag)
		le.Tags = appendIfMissing(le.Tags, ftag)
	}
}

// MarshalJSON Marshal LogEvent to Json
func (le LogEvent) MarshalJSON() (data []byte, err error) {
	event := le.getJSONMap()
	return json.Marshal(event)
}

// MarshalIndent Marshal LogEvent to Indent
func (le LogEvent) MarshalIndent() (data []byte, err error) {
	event := le.getJSONMap()
	return json.MarshalIndent(event, "", "\t")
}

// Get Value form LogEvent'Key
func (le LogEvent) Get(field string) (v interface{}) {
	switch field {
	case "@timestamp":
		v = le.Timestamp
	case "message":
		v = le.Message
	default:
		v = le.Extra[field]
	}
	return
}

// GetString Get Value-String form LogEvent'Key
func (le LogEvent) GetString(field string) (v string) {
	switch field {
	case "@timestamp":
		v = le.Timestamp.UTC().Format(timeFormat)
	case "message":
		v = le.Message
	default:
		if value, ok := le.Extra[field]; ok {
			v = fmt.Sprintf("%v", value)
		}
	}
	return
}

func appendIfMissing(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

func (le LogEvent) getJSONMap() map[string]interface{} {
	event := map[string]interface{}{
		"@timestamp": le.Timestamp.UTC().Format(timeFormat),
	}
	if le.Message != "" {
		event["message"] = le.Message
	}
	if len(le.Tags) > 0 {
		event["tags"] = le.Tags
	}
	for key, value := range le.Extra {
		event[key] = value
	}
	return event
}

// FormatWithEnv format string with environment value, ex: %{HOSTNAME}
func FormatWithEnv(text string) (result string) {
	result = text

	matches := revar.FindAllStringSubmatch(result, -1)
	for _, submatches := range matches {
		field := submatches[1]
		value := os.Getenv(field)
		if len(value) > 0 {
			result = strings.Replace(result, submatches[0], value, -1)
		}
	}

	return
}

// FormatWithTime format string with current time, ex: %{+2006-01-02}
func FormatWithTime(text string) (result string) {
	result = text

	matches := retime.FindAllStringSubmatch(result, -1)
	for _, submatches := range matches {
		value := time.Now().Format(submatches[1])
		result = strings.Replace(result, submatches[0], value, -1)
	}

	return
}

// Format return string with current time / LogEvent field / ENV, ex: %{hostname}
func (le LogEvent) Format(format string) (out string) {
	out = format

	out = FormatWithTime(out)

	matches := revar.FindAllStringSubmatch(out, -1)
	for _, submatches := range matches {
		field := submatches[1]
		value := le.GetString(field)
		if len(value) > 0 {
			out = strings.Replace(out, submatches[0], value, -1)
		}
	}

	out = FormatWithEnv(out)

	return
}
