package schema

import (
	"fmt"
	"strings"
	"time"

	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

// BusUnit represents the active bus unit of a bus company and the
// capacity for the specific unit.
//
// The "dynamodbav" struct tag can be used to control the value
// that will be marshaled into a AttributeValue.
type BusUnit struct {
	ID          string `json:"id" dynamodbav:"id"`                     // Unique bus unit ID as the primary key
	BusID       string `json:"bus_id" dynamodbav:"bus"`                // The Bus ID as the sort key
	Code        string `json:"code" dynamodbav:"code"`                 // Code is a uniqe identification of a bus unit
	Active      *bool  `json:"active" dynamodbav:"active"`             // Whether the bus unit is on trip and accepts a true or false value
	Capacity    int    `json:"capacity" dynamodbav:"capacity"`         // The number of passenger of a bus unit
	DateCreated string `json:"date_created" dynamodbav:"date_created"` // The date it was created as unix epoch time
}

func (unit BusUnit) Error(err error, code, message string, kv ...utility.KVP) {
	kv = append(kv, utility.KVP{Key: "Integration", Value: "Bus Ticketing – Bus Unit"})
	utility.Error(err, code, message, kv...)
}

// SetValues automatically generates the Bus Unit ID as your primary
// key and set the date it was created as unix epoch time.
//
// Example:
//		code: xyz-bus-0001
//		bus: BCDFGH-587390
//		id: BCDFGH-XYZ-BUS-0001
//		date_created: 1658837116
func (unit *BusUnit) SetValues() {
	unit.DateCreated = fmt.Sprint(time.Now().Unix())
	unit.ID = fmt.Sprintf("%s%s", strings.ToUpper(unit.Code), unit.DateCreated[2:8])
}
