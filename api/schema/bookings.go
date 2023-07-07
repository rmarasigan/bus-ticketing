package schema

import (
	"errors"
	"fmt"
	"strings"
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
	ID            string           `json:"id" dynamodbav:"id"`                                                 // Unique booking ID as the primary key
	UserID        string           `json:"user_id" dynamodbav:"user_id"`                                       // The user ID
	BusID         string           `json:"bus_id" dynamodbav:"bus_id"`                                         // The unique Bus ID
	BusRouteID    string           `json:"bus_route_id" dynamodbav:"bus_route_id"`                             // The unique Bus Route ID as the sort key
	Status        BookingStatus    `json:"status" dynamodbav:"status"`                                         // The status of the particular booking
	SeatNumber    string           `json:"seat_number" dynamodbav:"seat_number"`                               // The specific seat number(s) for the particular booking
	TravelDate    string           `json:"travel_date" dynamodbav:"travel_date"`                               // The date when to travel
	DateCreated   string           `json:"date_created" dynamodbav:"date_created"`                             // The date it was created as unix epoch time
	DateConfirmed string           `json:"date_confirmed,omitempty" dynamodbav:"date_confirmed,omitemptyelem"` // The date the booking was confirmed
	IsCancelled   *bool            `json:"is_cancelled,omitempty" dynamodbav:"is_cancelled,omitemptyelem"`     // Indicates if the booking is cancelled or not
	Cancelled     BookingCancelled `json:"cancelled,omitempty" dynamodbav:"-"`                                 // Contains the cancelled booking record
	Timestamp     string           `json:"timestamp" dynamodbav:"timestamp"`                                   // The timestamp when the request was made
}

// Cancelled contains the cancelled booking information.
type BookingCancelled struct {
	ID            string `json:"id" dynamodbav:"id"`                         // Unique booking cancellation ID
	BookingID     string `json:"booking_id" dynamodbav:"booking_id"`         // Unique booking ID as the primary key
	Reason        string `json:"reason" dynamodbav:"reason"`                 // Reason for booking cancellation
	CancelledBy   string `json:"cancelled_by" dynamodbav:"cancelled_by"`     // Indicates who cancelled the booking
	DateCancelled string `json:"date_cancelled" dynamodbav:"date_cancelled"` // The date when the booking was cancelled
}

// Error sets the default key-value pair.
func (booking Bookings) Error(err error, code, message string, kv ...utility.KVP) {
	if booking != (Bookings{}) {
		kv = append(kv, utility.KVP{Key: "bookings", Value: booking})
	}

	kv = append(kv, utility.KVP{Key: "Integration", Value: "Bus Ticketing – Bookings"})
	utility.Error(err, code, message, kv...)
}

// IsEmptyPayload checks if the request payload is empty and if it is,
// it will return an error message.
func (booking Bookings) IsEmptyPayload(payload string) error {
	if payload == "" {
		err := errors.New("payload is required")
		booking.Error(err, "APIError", "the reuqest payload is empty")

		return err
	}

	return nil
}

// IsValidStatus validates if the booking status is valid or not. If it
// is an invalid booking status, it will return an error message.
//
// Valid Booking Status:
//  PENDING, CONFIRMED, CANCELLED
func (booking Bookings) IsValidStatus() error {
	switch booking.Status {
	case booking.Status.Pending(),
		booking.Status.Confirmed(),
		booking.Status.Cancelled():

		return nil

	default:
		return errors.New("invalid booking status")
	}
}

// IsBookingCancelled validates if the details for the canceled booking
// are set or not. If the required fields are not set, it will return an
// error message.
//
//  Required fields:
//   reason, cancelled_by
func (booking Bookings) IsBookingCancelled() error {
	if booking.IsCancelled != nil && *booking.IsCancelled {
		// Check if it is set in the request payload
		if booking.Cancelled == (BookingCancelled{}) {
			return errors.New("'cancelled' fields are not set in the payload")

		} else {
			var msg []string

			if booking.Cancelled.Reason == "" {
				msg = append(msg, "'reason'")
			}

			if booking.Cancelled.CancelledBy == "" {
				msg = append(msg, "'cancelled_by'")
			}

			if len(msg) > 0 {
				return fmt.Errorf("object has missing required properties: [%s]", strings.Join(msg, ", "))
			}

			return nil
		}
	}

	return nil
}

// EventSource returns an EventBridge event source.
//
//  EventSource Types:
//   - booking:confirmed
//   - booking:cancelled
func (booking Bookings) EventSource() (string, error) {
	switch booking.Status {
	case booking.Status.Confirmed():
		return "booking:confirmed", nil

	case booking.Status.Cancelled():
		return "booking:cancelled", nil

	default:
		return "", fmt.Errorf("invalid booking event source [valid: %s, %s]", booking.Status.Confirmed(), booking.Status.Cancelled())
	}
}

// SetValues automatically generates the Bookings ID as your primary
// key, and set the date it was created as unix epoch time.
//
// Example:
//  ID: 36bc8bc5-44d6-447b-a63b-039b99658b78
//  DateCreated: 1688091891
func (booking *Bookings) SetValues() {
	booking.ID = uuid.NewString()
	booking.DateCreated = time.Now().Format("2006-01-02 15:04:05")
}
