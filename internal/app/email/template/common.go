package template

import (
	"fmt"

	"github.com/rmarasigan/bus-ticketing/api/schema"
)

// bookingDetails sets and returns the common email content details. The common details
// are user, route, and booking.
func bookingDetails(user schema.User, route schema.BusRoute, booking schema.Bookings) string {
	var details string

	details += fmt.Sprintf("<b>Passenger Name</b>: %s %s\n", user.FirstName, user.LastName)
	details += fmt.Sprintf("<b>Bus Number</b>: %s\n", route.BusUnitID)
	details += fmt.Sprintf("<b>Seat Number(s)</b>: %s\n\n", booking.SeatNumber)

	// *********************************************************** //
	// ******************** Departure Detials ******************** //
	// *********************************************************** //
	details += "<b>Departure Details</b>\n"
	details += fmt.Sprintf("\t\tLocation: %s\n", route.FromRoute)
	details += fmt.Sprintf("\t\tTime:&nbsp;&nbsp;\t%s\n\n", route.DepartureTime)

	// *********************************************************** //
	// ********************* Arrival Detials ********************* //
	// *********************************************************** //
	details += "<b>Arrival Details</b>\n"
	details += fmt.Sprintf("\t\tLocation: %s\n", route.ToRoute)
	details += fmt.Sprintf("\t\tTime:&nbsp;&nbsp;\t%s\n\n", route.ArrivalTime)

	return details
}
