package utility

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	LOG_OK    = "OK"
	LOG_INFO  = "INFO"
	LOG_DEBUG = "DEBUG"
	LOG_ERROR = "ERROR"
)

type Logs struct {
	Code         string         `json:"log_code"`
	Message      any            `json:"log_msg"`
	ErrorMessage string         `json:"log_errmsg,omitempty"`
	Level        string         `json:"log_level"`
	Keys         map[string]any `json:"log_keys,omitempty"`
	TimeStamp    string         `json:"log_timestamp"`
}

// Print marshal response JSON to print a string format JSON.
func (log *Logs) Print() {
	encodeJSON, err := json.Marshal(log)
	if err != nil {
		fmt.Println("Logger Print function failed to encode JSON")
	}

	fmt.Println(string(encodeJSON))
}

// SetKeys checks if Log Keys are empty in order to create an empty map. If it's not empty, set its key-value pair.
func (l *Logs) SetKeys(key string, value any) {
	if l.Keys == nil {
		// Create an empty map
		l.Keys = make(map[string]any)
	}

	// Set key-value pairs using typical name[key] = val syntax
	l.Keys[key] = value
}

// SetTimeStamp sets the current timestamp with the ff. format: 2006-01-02 15:04:05.
func (l *Logs) SetTimeStamp() {
	l.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
}

// OK prints an OK log information
func OK(code, message string, kv ...KVP) {
	var entry Logs

	entry.Code = code
	entry.Level = LOG_OK
	entry.Message = message
	entry.SetTimeStamp()

	if len(kv) != 0 {
		for _, kvp := range kv {
			entry.SetKeys(kvp.KeyValue())
		}
	}

	entry.Print()
}

// Info prints a log information.
func Info(code, message string, kv ...KVP) {
	var entry Logs

	entry.Code = code
	entry.Level = LOG_INFO
	entry.Message = message
	entry.SetTimeStamp()

	if len(kv) != 0 {
		for _, kvp := range kv {
			entry.SetKeys(kvp.KeyValue())
		}
	}

	entry.Print()
}

// Debug prints a debug log information.
func Debug(code, message string, kv ...KVP) {
	var entry Logs

	entry.Code = code
	entry.Level = LOG_DEBUG
	entry.Message = message
	entry.SetTimeStamp()

	if len(kv) != 0 {
		for _, kvp := range kv {
			entry.SetKeys(kvp.KeyValue())
		}
	}

	entry.Print()
}

// Error prints an error log information.
func Error(err error, code, message string, kv ...KVP) {
	var entry Logs

	entry.Level = LOG_ERROR
	entry.Code = code
	entry.Message = message
	entry.SetTimeStamp()

	if err != nil {
		entry.ErrorMessage = err.Error()
	}

	if len(kv) != 0 {
		for _, kvp := range kv {
			entry.SetKeys(kvp.KeyValue())
		}
	}

	entry.Print()
}
