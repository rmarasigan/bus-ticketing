package schema

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rmarasigan/bus-ticketing/internal/app"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

// BusRoute is used to store the specific bus unit route, rate, and schedule.
//
// The "dynamodbav" struct tag can be used to control the value
// that will be marshaled into a AttributeValue.
type BusRoute struct {
	ID            string   `json:"id" dynamodbav:"id"`                                             // Unique bus route ID as the primary key
	BusID         string   `json:"bus_id" dynamodbav:"bus_id"`                                     // The Bus ID as the sort key
	BusUnitID     string   `json:"bus_unit_id" dynamodbav:"bus_unit_id"`                           // The Bus Unit ID for the identification of specific bus unit route
	Currency      string   `json:"currency_code" dynamodbav:"currency_code"`                       // Medium of exchange for goods and services
	Rate          *float64 `json:"rate" dynamodbav:"rate"`                                         // Fare charged to the passenger
	Active        *bool    `json:"active" dynamodbav:"active"`                                     // Defines if the bus is available for that route
	DepartureTime string   `json:"departure_time" dynamodbav:"departure_time"`                     // Expected departure time on the starting point and in 24-hour format
	ArrivalTime   string   `json:"arrival_time" dynamodbav:"arrival_time"`                         // Expected arrival time on the destination and in 24-hour format
	FromRoute     string   `json:"from_route" dynamodbav:"from_route"`                             // Indicating the starting point of a bus
	ToRoute       string   `json:"to_route" dynamodbav:"to_route"`                                 // Indicating the destination of bus
	DateCreated   string   `json:"date_created,omitempty" dynamodbav:"date_created,omitemptyelem"` // The date it was created as unix epoch time
}

// Error sets the default key-value pair.
func (route BusRoute) Error(err error, code, message string, kv ...utility.KVP) {
	if route != (BusRoute{}) {
		kv = append(kv, utility.KVP{Key: "bus_route", Value: route})
	}

	kv = append(kv, utility.KVP{Key: "Integration", Value: "Bus Ticketing – Bus Route"})
	utility.Error(err, code, message, kv...)
}

// IsEmptyPayload checks if the request payload is empty and if it is,
// it will return an error message.
func (route BusRoute) IsEmptyPayload(payload string) error {
	if payload == "" {
		err := errors.New("payload is required")
		route.Error(err, "APIError", "the request payload is empty")

		return err
	}

	return nil
}

// primaryKey uses from_route, to_route, departure_time and arrival_time
// to form the Bus Route key.
//
// Example:
//		id: RTRTB15001900877732
func (route BusRoute) primaryKey() string {
	var key string

	from, err := app.RemoveVowel(route.FromRoute)
	if err != nil {
		route.Error(err, "Key", "failed to remove vowel letters.")
		return ""
	}

	from, err = app.RemoveSymbols(from)
	if err != nil {
		route.Error(err, "Key", "failed to remove symbols.")
		return ""
	}

	to, err := app.RemoveVowel(route.ToRoute)
	if err != nil {
		route.Error(err, "Key", "failed to remove vowel letters.")
		return ""
	}

	to, err = app.RemoveSymbols(to)
	if err != nil {
		route.Error(err, "Key", "failed to remove symbols.")
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
	route.ID = route.primaryKey()
}

// BusRouteFilter contains the fields of a bus route
// that can be used for filtering.
type BusRouteFilter struct {
	BusID     string
	BusUnitID string
	Active    *bool
	Departure string
	Arrival   string
	FromRoute string
	ToRoute   string
}

// SetFilter sets and returns the fields of the bus route
// that will be used in filtering data.
func (route BusRoute) SetFilter() BusRouteFilter {
	return BusRouteFilter{
		BusID:     route.BusID,
		BusUnitID: route.BusUnitID,
		Active:    route.Active,
		Departure: route.DepartureTime,
		Arrival:   route.ArrivalTime,
		FromRoute: route.FromRoute,
		ToRoute:   route.ToRoute,
	}
}
