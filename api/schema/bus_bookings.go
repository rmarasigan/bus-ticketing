package schema

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

// BookingStatus is the status of the particular booking.
type BookingStatus string

// Pending booking status means that it is not yet
// confirmed of processed.
func (BookingStatus) Pending() BookingStatus {
	return "PENDING"
}

// WaitingForPayment booking status means that it is still
// waiting for the payment confirmation.
func (BookingStatus) WaitingForPayment() BookingStatus {
	return "WAITING_FOR_PAYMENT"
}

// Confirmed booking status means that is confirmed by
// the administrator.
func (BookingStatus) Confirmed() BookingStatus {
	return "CONFIRMED"
}

// Cancelled booking status means that the booking is cancelled
// by the user or due to other reason.
func (BookingStatus) Cancelled() BookingStatus {
	return "CANCELLED"
}

// Bookings is used to store the details of reserving seats for a
// particular bus.
//
// The "dynamodbav" struct tag can be used to control the value that
// will be marshaled into a AttributeValue.
type Bookings struct {
	ID          string        `json:"id" dynamodbav:"id"`                                 // Unique booking ID as the primary key
	UserID      string        `json:"user_id" dynamodbav:"user_id"`                       // The user ID
	BusID       string        `json:"bus_id" dynamodbav:"bus_id"`                         // The unique Bus ID
	BusRouteID  string        `json:"bus_route_id" dynamodbav:"bus_route_id"`             // The unique Bus Route ID as the sort key
	Status      BookingStatus `json:"status" dynamodbav:"status"`                         // The status of the particular booking
	SeatNumber  string        `json:"seat_number" dynamodbav:"seat_number"`               // The specific seat number(s) for the particular booking
	Reason      string        `json:"reason,omitempty" dynamodbav:"reason,omitemptyelem"` // The reason why the booking is cancelled
	Timestamp   string        `json:"timestamp" dynamodbav:"timestamp"`                   // The timestamp when the request was made
	DateCreated string        `json:"date_created" dynamodbav:"date_created"`             // The date it was created as unix epoch time
}

// Error sets the default key-value pair.
func (booking Bookings) Error(err error, code, message string, kv ...utility.KVP) {
	if booking != (Bookings{}) {
		kv = append(kv, utility.KVP{Key: "bookings", Value: booking})
	}

	kv = append(kv, utility.KVP{Key: "Integration", Value: "Bus Ticketing – Bookings"})
	utility.Error(err, code, message, kv...)
}

// SetValues automatically generates the Bookings ID as your primary
// key, and set the date it was created as unix epoch time.
//
// Example:
//  ID: 36bc8bc5-44d6-447b-a63b-039b99658b78
//  DateCreated: 1688091891
func (booking *Bookings) SetValues() {
	booking.ID = uuid.NewString()
	booking.DateCreated = fmt.Sprint(time.Now().Unix())
}
