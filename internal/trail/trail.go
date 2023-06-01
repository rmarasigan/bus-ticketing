package trail

import (
	"fmt"
	"time"

	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

const (
	OK    = 0
	INFO  = 1
	DEBUG = 2
	ERROR = 3
)

var trailLevel = map[int]string{
	OK:    "OK",
	INFO:  "INFO",
	DEBUG: "DEBUG",
	ERROR: "ERROR",
}

type Trail struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	TimeStamp string `json:"timestamp"`
}

// SetTimeStamp sets the current timestamp with the ff. format: 2006-01-02 15:04:05.
func (t *Trail) SetTimeStamp() {
	t.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
}

// Print accepts a level parameter and formats according to a format.
//
// level accepts OK, INFO, DEBUG and ERROR.
func Print(level int, msg any, i ...any) {
	message := fmt.Sprint(msg)

	entry := new(Trail)
	entry.Level = trailLevel[level]
	entry.Message = fmt.Sprintf(message, i...)
	entry.SetTimeStamp()

	fmt.Println(utility.EncodeJSON(entry))
}

// Ok prints an OK information and formats the message.
func Ok(msg any, i ...any) {
	message := fmt.Sprint(msg)

	entry := new(Trail)
	entry.Level = trailLevel[OK]
	entry.Message = fmt.Sprintf(message, i...)
	entry.SetTimeStamp()

	fmt.Println(utility.EncodeJSON(entry))
}

// Info prints an information and formats the message.
func Info(msg any, i ...any) {
	message := fmt.Sprint(msg)

	entry := new(Trail)
	entry.Level = trailLevel[INFO]
	entry.Message = fmt.Sprintf(message, i...)
	entry.SetTimeStamp()

	fmt.Println(utility.EncodeJSON(entry))
}

// Debug prints a debug information and formats the message.
func Debug(msg any, i ...any) {
	message := fmt.Sprint(msg)

	entry := new(Trail)
	entry.Level = trailLevel[DEBUG]
	entry.Message = fmt.Sprintf(message, i...)
	entry.SetTimeStamp()

	fmt.Println(utility.EncodeJSON(entry))
}

// Error prints an error information and formats the message.
func Error(msg any, i ...any) {
	message := fmt.Sprint(msg)

	entry := new(Trail)
	entry.Level = trailLevel[ERROR]
	entry.Message = fmt.Sprintf(message, i...)
	entry.SetTimeStamp()

	fmt.Println(utility.EncodeJSON(entry))
}
