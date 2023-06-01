package schema

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

// BusTrip defines the current count of the capacity of a bus for the client to know
// if there are enough seats for them.
//
// The "dynamodbav" struct tag can be used to control the value
// that will be marshaled into a AttributeValue.
type BusTrip struct {
	ID          string `json:"id" dynamodbav:"id"`                     // Unique bus trip ID as the primary key
	BusUnit     string `json:"bus_unit_id" dynamodbav:"bus_unit_id"`   // The Bus Unit ID as sort key and for the identification of specific bus trip
	BusRoute    string `json:"bus_route_id" dynamodbav:"bus_route_id"` // The Bus Route ID
	SeatsLeft   int    `json:"seats_left" dynamodbav:"seats_left"`     // Current count left of seats
	DateCreated string `json:"date_created" dynamodbav:"date_created"` // The date it was created as unix epoch time
}

func (trip BusTrip) Error(err error, code, message string, kv ...utility.KVP) {
	kv = append(kv, utility.KVP{Key: "Integration", Value: "Bus Ticketing – Bus Trip"})
	utility.Error(err, code, message, kv...)
}

// SetValues automatically generates the Bus Trip ID as your primary
// key, set the date it was created as unix epoch time, and set the
// seats remaining.
//
// Function Parameters:
//  capacity: The number of bus capacity that a bus can accomodate
//  seat: The requested number of seats to be reserved by the customer
func (trip *BusTrip) SetValues(capacity, seat int) {
	trip.SeatsLeft = (capacity - seat)
	trip.DateCreated = fmt.Sprint(time.Now().Unix())
	trip.ID = fmt.Sprintf("%s%s", trip.BusUnit[(len(trip.BusUnit)-6):len(trip.BusUnit)], trip.DateCreated[3:9])
}

// SetDefaultValues automatically sets the old or default value for the
// ID, DateCreated and SeatsLeft.
//
// Function Parameters:
//  old: The Bus Trip information that is within the day
//  seat: The requested number of seats to be reserved by the customer
func (trip *BusTrip) SetDefaultValues(seat int, old *BusTrip) error {
	trip.ID = old.ID
	trip.DateCreated = old.DateCreated

	if old.SeatsLeft <= 0 || seat > old.SeatsLeft {
		return errors.New("there are not enough seats available")
	}

	if old.SeatsLeft > 0 && seat <= old.SeatsLeft {
		trip.SeatsLeft = (old.SeatsLeft - seat)
	}

	return nil
}

// IsWithinTheDay validate whether the date and time the trip was created is within the
// current date range.
func (trip *BusTrip) IsWithinTheDay() (bool, error) {
	now := time.Now()
	layout := "3:04 PM"

	startTime, err := time.Parse(layout, "12:00 AM")
	if err != nil {
		trip.Error(err, "IsWithinTheDay", "cannot parse string time to time.Time.")
		return false, err
	}

	endTime, err := time.Parse(layout, "11:59 PM")
	if err != nil {
		trip.Error(err, "IsWithinTheDay", "cannot parse string to time.Time.")
		return false, err
	}

	ending := time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), 0, 0, now.Location())
	beginning := time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, now.Location())

	created, err := strconv.Atoi(trip.DateCreated)
	if err != nil {
		trip.Error(err, "IsWithinTheDay", "failed to convert date created string to int.", utility.KVP{Key: "date_created", Value: trip.DateCreated})
		return false, err
	}
	dateCreated := time.Unix(int64(created), 0)

	return (dateCreated.After(beginning) && dateCreated.Before(ending)), nil
}
