package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rmarasigan/bus-ticketing/pkg/cw/kvp"
)

const (
	LOG_INFO     = "INFO"
	LOG_DEBUG    = "DEBUG"
	LOG_ERROR    = "ERROR"
	LOG_WARNING  = "WARNING"
	LOG_CRITICAL = "CRITICAL"
)

type Logs struct {
	Code         string                 `json:"log_code"`
	Message      interface{}            `json:"log_msg"`
	ErrorMessage string                 `json:"log_errmsg,omitempty"`
	Level        string                 `json:"log_level"`
	Keys         map[string]interface{} `json:"log_kvp,omitempty"`
	TimeStamp    string                 `json:"log_timestamp"`
}

// print marshal a Log struct to print a string format JSON.
func (logs *Logs) print() {
	encodeJSON, err := json.Marshal(logs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(encodeJSON))
}

// setKeys checks if Log Keys are empty in order to create an empty map. If it's not empty, set its key-value pair.
func (l *Logs) setKeys(key string, value interface{}) {
	if l.Keys == nil {
		// Create an empty map
		l.Keys = make(map[string]interface{})
	}

	// Set key-value pairs using typical name[key] = val syntax
	l.Keys[key] = value
}

// setTimeStamp sets the current timestamp with the ff. format 2006-01-02 15:04:05.
func (l *Logs) setTimeStamp() {
	l.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
}

// Info logs an information.
func Info(logs *Logs, kv ...kvp.Attribute) {
	var entry Logs

	entry.Code = logs.Code
	entry.Level = LOG_INFO
	entry.Message = logs.Message
	entry.setTimeStamp()

	if len(kv) != 0 {
		for _, kvp := range kv {
			entry.setKeys(kvp.KeyValue())
		}
	}

	entry.print()
}

// Debug logs a debug information.
func Debug(logs *Logs, kv ...kvp.Attribute) {
	var entry Logs

	entry.Code = logs.Code
	entry.Level = LOG_DEBUG
	entry.Message = logs.Message
	entry.setTimeStamp()

	if len(kv) != 0 {
		for _, kvp := range kv {
			entry.setKeys(kvp.KeyValue())
		}
	}

	entry.print()
}

// Error logs an error information.
func Error(err error, logs *Logs, kv ...kvp.Attribute) {
	var entry Logs

	entry.Code = logs.Code
	entry.Level = LOG_ERROR
	entry.Message = logs.Message
	entry.setTimeStamp()

	if err != nil {
		entry.ErrorMessage = err.Error()
	}

	if len(kv) != 0 {
		for _, kvp := range kv {
			entry.setKeys(kvp.KeyValue())
		}
	}

	entry.print()
}

// Warning logs a warning information.
func Warning(logs *Logs, kv ...kvp.Attribute) {
	var entry Logs

	entry.Code = logs.Code
	entry.Level = LOG_WARNING
	entry.Message = logs.Message
	entry.setTimeStamp()

	if len(kv) != 0 {
		for _, kvp := range kv {
			entry.setKeys(kvp.KeyValue())
		}
	}

	entry.print()
}

// Critical logs a critical information.
func Critical(logs *Logs, kv ...kvp.Attribute) {
	var entry Logs

	entry.Code = logs.Code
	entry.Level = LOG_CRITICAL
	entry.Message = logs.Message
	entry.setTimeStamp()

	if len(kv) != 0 {
		for _, kvp := range kv {
			entry.setKeys(kvp.KeyValue())
		}
	}

	entry.print()
}
