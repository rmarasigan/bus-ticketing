package schema

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rmarasigan/bus-ticketing/internal/app"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

// Bus contains the basic information of the bus company.
//
// The "dynamodbav" struct tag can be used to control the value
// that will be marshaled into a AttributeValue.
type Bus struct {
	ID           string `json:"id" dynamodbav:"id"`                                             // Unique bus ID
	Name         string `json:"name" dynamodbav:"name"`                                         // Name of the bus line as the primary key
	Owner        string `json:"owner" dynamodbav:"owner"`                                       // Bus company owner and is a required field
	Email        string `json:"email" dynamodbav:"email"`                                       // Bus company email and is a required field
	Address      string `json:"address" dynamodbav:"address"`                                   // Bus company address and is a required field
	Company      string `json:"company" dynamodbav:"company"`                                   // Name of the company and serves as your sort key and is a required field
	MobileNumber string `json:"mobile_number" dynamodbav:"mobile_number"`                       // Bus company mobile number and is a required field
	DateCreated  string `json:"date_created,omitempty" dynamodbav:"date_created,omitemptyelem"` // The date it was created as unix epoch time
}

// Error sets the default key-value pair.
func (bus Bus) Error(err error, code, message string, kv ...utility.KVP) {
	if bus != (Bus{}) {
		kv = append(kv, utility.KVP{Key: "bus", Value: bus})
	}

	kv = append(kv, utility.KVP{Key: "Integration", Value: "Bus Ticketing – Bus"})
	utility.Error(err, code, message, kv...)
}

// IsEmptyPayload checks if the request payload is empty and if it is,
// it will return an error message.
func (bus Bus) IsEmptyPayload(payload string) error {
	if payload == "" {
		err := errors.New("payload is required")
		bus.Error(err, "APIError", "the request payload is empty")

		return err
	}

	return nil
}

// partialPrimaryKey uses company name, removes the vowel letters and
// converts it to uppercase to generate a key value of bus that will
// be used for the bus ID.
//
// Example:
//		id: RLBSW-856996
func (bus Bus) partialPrimaryKey() string {
	key, err := app.RemoveVowel(bus.Company)
	if err != nil {
		bus.Error(err, "Key", "failed to remove vowel letters.")
		return ""
	}

	key, err = app.RemoveSymbols(key)
	if err != nil {
		bus.Error(err, "Key", "failed to remove symbols.")
		return ""
	}

	return strings.ToUpper(key)
}

// SetValues automatically generates the Bus ID as your primary key,
// and set the date it was created as unix epoch time.
//
// Example:
//		id: RLBSW-856996
//		date_created: 1685699666
func (bus *Bus) SetValues() {
	bus.DateCreated = fmt.Sprint(time.Now().Unix())
	bus.ID = fmt.Sprintf("%s-%s", bus.partialPrimaryKey(), bus.DateCreated[2:8])
}

// FailedBus represents the failed bus that needs to be re-processed.
type FailedBus struct {
	Failed []struct {
		Bus    Bus    `json:"bus"`
		Reason string `json:"reason,omitempty"`
	} `json:"failed"`
}

// SetFailedBus sets the failed bus information transaction by passing the
// Bus and the reason why it failed.
func (failed *FailedBus) SetFailedBus(bus Bus, reason string) {
	data := struct {
		Bus    Bus    `json:"bus"`
		Reason string `json:"reason,omitempty"`
	}{}

	data.Bus = bus
	data.Reason = reason
	failed.Failed = append(failed.Failed, data)
}
