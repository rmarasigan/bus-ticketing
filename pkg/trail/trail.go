package trail

import (
	"fmt"
	"time"

	"github.com/rmarasigan/bus-ticketing/pkg/api"
)

const (
	OK       = 0
	INFO     = 1
	DEBUG    = 2
	WARNING  = 3
	ERROR    = 4
	CRITICAL = 5
)

var trailLevel = map[int]string{
	OK:       "OK",
	INFO:     "INFO",
	DEBUG:    "DEBUG",
	WARNING:  "WARNING",
	ERROR:    "ERROR",
	CRITICAL: "CRITICAL",
}

type Trail struct {
	Level     string `json:"trail_level"`
	Message   string `json:"trail_message"`
	TimeStamp string `json:"timestamp"`
}

// setTimeStamp sets the current timestamp with the ff. format: 2006-01-02 15:04:05.
func (t *Trail) setTimeStamp() {
	t.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
}

// Print accepts a level parameter and formats according to a format specifier.
//
// level accepts OK, INFO, DEBUG, WARNING, ERROR, CRITICAL.
func Print(level int, msg interface{}, i ...interface{}) {
	message := fmt.Sprint(msg)

	entry := new(Trail)
	entry.Level = trailLevel[level]
	entry.Message = fmt.Sprintf(message, i...)
	entry.setTimeStamp()

	fmt.Println(string(api.EncodeResponse(entry)))
}
