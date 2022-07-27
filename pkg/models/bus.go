package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/rmarasigan/bus-ticketing/pkg/common"
	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
)

// Bus contains the information of the bus company.
type Bus struct {
	ID           string `json:"id"`            // Unique bus ID as the primary key
	Owner        string `json:"owner"`         // Bus company owner and is a required field
	Email        string `json:"email"`         // Bus company email and is a required field
	Address      string `json:"address"`       // Bus company address and is a required field
	Company      string `json:"company"`       // Name of the company and serves as your sort key and is a required field
	MobileNumber string `json:"mobile_number"` // Bus company mobile number and is a required field
	DateCreated  string `json:"date_created"`  // The date it was created as unix epoch time
}

// Key uses company name, removes the vowel letters and converts
// it to uppercase to generate a key value of bus that will be used
// for the bus ID.
//
// Example:
//		company: abcdefghi
//		key: BCDFGH
func (bus Bus) Key() string {
	key, err := common.RemoveVowel(bus.Company)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "SetValues", Message: "Failed to remove vowel letters"})
		return ""
	}

	key = strings.Replace(strings.ToUpper(key), " ", "", -1)

	return key
}

// SetValues automatically generates the Bus ID as your primary key,
// and set the date it was created as unix epoch time.
//
// Example:
//		id: BCDFGH-587390
//		date_created: 1658739080
func (bus *Bus) SetValues() {
	bus.DateCreated = fmt.Sprint(time.Now().Unix())
	bus.ID = fmt.Sprintf("%s-%s", bus.Key(), bus.DateCreated[2:8])
}

// ValidateUpdate validates if the field that are going to be updated
// are empty or not to set its previous value.
//
// Fields that are validated:
//     * owner, address, email, mobile_number
func (bus *Bus) ValidateUpdate(old *Bus) {
	if bus.Owner == "" {
		bus.Owner = old.Owner
	}

	if bus.Address == "" {
		bus.Address = old.Address
	}

	if bus.Email == "" {
		bus.Email = old.Email
	}

	if bus.MobileNumber == "" {
		bus.MobileNumber = old.MobileNumber
	}
}

// BusUnit represents the active bus unit of a bus company and the
// capacity for the specific unit.
type BusUnit struct {
	ID          string `json:"id"`           // Unique bus unit ID as the primary key
	Bus         string `json:"bus"`          // The Bus ID as the sort key
	Code        string `json:"code"`         // Code is a uniqe identification of a bus unit
	Active      *bool  `json:"active"`       // Whether the bus unit is on trip and accepts a true or false value
	Capacity    int    `json:"capacity"`     // The number of passenger of a bus unit
	DateCreated string `json:"date_created"` // The date it was created as unix epoch time
}

// SetValues automatically generates the Bus Unit ID as your primary
// key, set the bus info, and set the date it was created as unix epoch time.
//
// Example:
//		code: xyz-bus-0001
//		bus: BCDFGH-587390
//		id: BCDFGH-XYZ-BUS-0001
//		date_created: 1658837116
func (unit *BusUnit) SetValues() {
	key := strings.Split(unit.Bus, "-")[0]

	unit.DateCreated = fmt.Sprint(time.Now().Unix())
	unit.ID = fmt.Sprintf("%s-%s", key, strings.ToUpper(unit.Code))
}

// ValidateUpdate validates bus unit field if they are empty or not
// to set its previous value.
//
// Fields that are validated:
//     * active, capacity
func (unit *BusUnit) ValidateUpdate(old *BusUnit) {
	if unit.Active == nil {
		unit.Active = old.Active
	}

	if unit.Capacity == 0 {
		unit.Capacity = old.Capacity
	}
}

type BusRoute struct {
	ID            string  `json:"id"`
	BusUnit       BusUnit `json:"bus_unit,omitempty"`
	Name          string  `json:"route_name"`
	Rate          float64 `json:"route_rate"`
	Available     bool    `json:"route_available"`
	DepartureTime string  `json:"route_departure_time"`
	ArrivalTime   string  `json:"route_arrival_time"`
	FromRoute     string  `json:"route_from"`
	ToRouteCode   string  `json:"route_to"`
	DateCreated   string  `json:"route_date_created"`
}

type BusTrip struct {
	ID          string   `json:"id"`
	BusRoute    BusRoute `json:"bus_route"`
	Status      string   `json:"trip_status"`
	SeatsLeft   int      `json:"trip_seats_left"`
	DateCreated string   `json:"trip_date_created"`
}
