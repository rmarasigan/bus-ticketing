package schema

import (
	"fmt"
	"time"

	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

const BUS_UNIT_MIN_CAPACITY = 25

// BusUnit represents a bus company's active bus unit and the specific
// unit's capacity. The "Active" is set to a boolean pointer for us to
// validate if it is set by checking if the value is "nil" since it will
// be the default value if it is uninitialized.
//
// Reference: https://stackoverflow.com/a/43351386/19679222
//
// The "dynamodbav" struct tag can be used to control the value
// that will be marshaled into a AttributeValue.
type BusUnit struct {
	BusID       string `json:"bus_id" dynamodbav:"bus_id"`                                     // The Bus ID as the sort key
	Code        string `json:"code" dynamodbav:"code"`                                         // Code is a uniqe identification of a bus unit
	Active      *bool  `json:"active" dynamodbav:"active"`                                     // Whether the bus unit is on trip and accepts a true or false value
	MinCapacity *int   `json:"min_capacity" dynamodbav:"min_capacity"`                         // The minimum number of passenger of a bus unit
	MaxCapacity *int   `json:"max_capacity" dynamodbav:"max_capacity"`                         // The maximum number of passenger of a bus unit
	DateCreated string `json:"date_created,omitempty" dynamodbav:"date_created,omitemptyelem"` // The date it was created as unix epoch time
}

// Error sets the default key-value pair.
func (unit BusUnit) Error(err error, code, message string, kv ...utility.KVP) {
	if unit != (BusUnit{}) {
		kv = append(kv, utility.KVP{Key: "bus_unit", Value: unit})
	}

	kv = append(kv, utility.KVP{Key: "Integration", Value: "Bus Ticketing – Bus Unit"})
	utility.Error(err, code, message, kv...)
}

// SetValues automatically generates the Bus Unit ID as your primary
// key and set the date it was created as unix epoch time.
//
// Example:
//		date_created: 1658837116
func (unit *BusUnit) SetValues() {
	unit.DateCreated = fmt.Sprint(time.Now().Unix())
}

// ValidateMinimumCapacity validates if the amount set for the minimum
// capacity is not a negative amount or less than the minimum capacity
// requirement of the Bus Unit.
func (unit BusUnit) ValidateMinimumCapacity() error {
	if unit.MinCapacity != nil && *unit.MinCapacity < BUS_UNIT_MIN_CAPACITY {
		return fmt.Errorf("cannot set %v as the minimum capacity", *unit.MinCapacity)
	}

	return nil
}

// ValidateMaximumCapacity validates if the amount set for the maximum
// capacity is not a negative amount or less than the existing minimum
// capacity.
func (unit BusUnit) ValidateMaximumCapacity(min_capacity int) error {
	if unit.MaxCapacity != nil && *unit.MaxCapacity < min_capacity {
		return fmt.Errorf("cannot set %v as the max capacity that is lower than the minimum capacity", *unit.MaxCapacity)
	}

	return nil
}

// FailedBusUnits represents the failed bus unit that needs to be re-processed.
type FailedBusUnits struct {
	Failed []struct {
		Unit   BusUnit `json:"unit"`
		Reason string  `json:"reason,omitempty"`
	} `json:"failed"`
}

// SetFailedUnits sets the failed bus units information transaction by passing the
// Bus Unit and the reason why it failed.
func (failed *FailedBusUnits) SetFailedUnits(unit BusUnit, reason string) {
	data := struct {
		Unit   BusUnit `json:"unit"`
		Reason string  `json:"reason,omitempty"`
	}{}

	data.Unit = unit
	data.Reason = reason
	failed.Failed = append(failed.Failed, data)
}
