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
func (bus Bus) Key() string {
	key, err := common.RemoveVowel(bus.Company)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "BusKey", Message: "Failed to remove vowel letters."})
		return ""
	}

	key, err = common.RemoveSymbols(key)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "BusKey", Message: "Failed to remove symbols."})
		return ""
	}

	key = strings.ToUpper(key)

	return key
}

// SetValues automatically generates the Bus ID as your primary key,
// and set the date it was created as unix epoch time.
func (bus *Bus) SetValues() {
	bus.DateCreated = fmt.Sprint(time.Now().Unix())
	bus.ID = fmt.Sprintf("%s-%s", bus.Key(), bus.DateCreated[2:8])
}

// ValidateUpdate validates if the field that are going to be updated
// are empty or not to set its previous value.
//
// Fields that are validated:
//  owner, address, email, mobile_number
func (bus *Bus) ValidateUpdate(old *Bus) {
	if bus.Company == "" {
		bus.Company = old.Company
	}

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
	unit.DateCreated = fmt.Sprint(time.Now().Unix())
	unit.ID = fmt.Sprintf("%s%s", strings.ToUpper(unit.Code), unit.DateCreated[2:8])
}

// ValidateUpdate validates bus unit field if they are empty or not
// to set its previous value.
//
// Fields that are validated:
//  active, capacity
func (unit *BusUnit) ValidateUpdate(old *BusUnit) {
	if unit.Bus == "" {
		unit.Bus = old.Bus
	}

	if unit.Code == "" {
		unit.Code = old.Code
	}

	if unit.Active == nil {
		unit.Active = old.Active
	}

	if unit.Capacity == 0 {
		unit.Capacity = old.Capacity
	}
}

// BusRoute is used to store the specific bus unit route, rate, and schedule.
type BusRoute struct {
	ID            string  `json:"id"`             // Unique bus route ID as the primary key
	Bus           string  `json:"bus"`            // The Bus ID as the sort key
	BusUnit       string  `json:"bus_unit"`       // The Bus Unit ID for the identification of specific bus unit route
	Currency      string  `json:"currency_code"`  // Medium of exchange for goods and services
	Rate          float64 `json:"rate"`           // Fare charged to the passenger
	Available     *bool   `json:"available"`      // Defines if the bus is available for that route
	DepartureTime string  `json:"departure_time"` // Expected departure time on the starting point
	ArrivalTime   string  `json:"arrival_time"`   // Expected arrival time on the destination
	FromRoute     string  `json:"from_route"`     // Indicating the starting point of a bus and in 24-hour format
	ToRoute       string  `json:"to_route"`       // Indicating the destination of bus and in 24-hour format
	DateCreated   string  `json:"date_created"`   // The date it was created as unix epoch time
}

func (route BusRoute) Key() string {
	var key string

	from, err := common.RemoveVowel(route.FromRoute)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "BusRouteKey", Message: "Failed to remove vowel letters."})
		return ""
	}

	from, err = common.RemoveSymbols(from)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "BusRouteKey", Message: "Failed to remove symbols."})
		return ""
	}

	to, err := common.RemoveVowel(route.ToRoute)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "BusRouteKey", Message: "Failed to remove vowel letters."})
		return ""
	}

	to, err = common.RemoveSymbols(to)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "BusRouteKey", Message: "Failed to remove symbols."})
		return ""
	}

	to = strings.ToUpper(to)
	from = strings.ToUpper(from)
	departure := strings.ReplaceAll(route.DepartureTime, ":", "")
	arrival := strings.ReplaceAll(route.ArrivalTime, ":", "")
	key = fmt.Sprintf("%s%s%s%s%s", from, to, departure, arrival, route.DateCreated[2:8])

	return key
}

// SetValues automatically generates the Bus Route ID as your primary
// key, and set the date it was created as unix epoch time.
func (route *BusRoute) SetValues() {
	route.DateCreated = fmt.Sprint(time.Now().Unix())
	route.ID = route.Key()
}

// ValidateUpdate validates bus route field if they are empty or not
// to set its previous value.
//
// Fields that are validated:
//  rate, currency_code, available, departure_time, arrival_time, from_route, to_route
func (route *BusRoute) ValidateUpdate(old *BusRoute) {
	if route.ID == "" {
		route.ID = old.ID
	}

	if route.Rate <= 0 {
		route.Rate = old.Rate
	}

	if route.Currency == "" {
		route.Currency = old.Currency
	}

	if route.Available == nil {
		route.Available = old.Available
	}

	if route.DepartureTime == "" {
		route.DepartureTime = old.DepartureTime
	}

	if route.ArrivalTime == "" {
		route.ArrivalTime = old.ArrivalTime
	}

	if route.FromRoute == "" {
		route.FromRoute = old.FromRoute
	}

	if route.ToRoute == "" {
		route.ToRoute = old.ToRoute
	}
}

type BusTrip struct {
	ID          string   `json:"id"`
	BusRoute    BusRoute `json:"bus_route"`
	Status      string   `json:"trip_status"`
	SeatsLeft   int      `json:"trip_seats_left"`
	DateCreated string   `json:"trip_date_created"`
}
