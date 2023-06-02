package schema

import (
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
	ID           string `json:"id" dynamodbav:"id"`                       // Unique bus ID as the primary key
	Name         string `json:"name" dynamodbav:"name"`                   // Name of the bus line
	Owner        string `json:"owner" dynamodbav:"owner"`                 // Bus company owner and is a required field
	Email        string `json:"email" dynamodbav:"email"`                 // Bus company email and is a required field
	Address      string `json:"address" dynamodbav:"address"`             // Bus company address and is a required field
	Company      string `json:"company" dynamodbav:"company"`             // Name of the company and serves as your sort key and is a required field
	MobileNumber string `json:"mobile_number" dynamodbav:"mobile_number"` // Bus company mobile number and is a required field
	DateCreated  string `json:"date_created" dynamodbav:"date_created"`   // The date it was created as unix epoch time
}

// Error sets the default key-value pair.
func (bus Bus) Error(err error, code, message string, kv ...utility.KVP) {
	if bus != (Bus{}) {
		kv = append(kv, utility.KVP{Key: "bus", Value: bus})
	}

	kv = append(kv, utility.KVP{Key: "Integration", Value: "Bus Ticketing – Bus"})
	utility.Error(err, code, message, kv...)
}

// key uses company name, removes the vowel letters and converts
// it to uppercase to generate a key value of bus that will be used
// for the bus ID.
//
// Example:
//		id: RLBSW-856996
func (bus Bus) key() string {
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
	bus.ID = fmt.Sprintf("%s-%s", bus.key(), bus.DateCreated[2:8])
}
