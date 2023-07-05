package template

import (
	"fmt"
	"strings"

	"github.com/rmarasigan/bus-ticketing/api/schema"
)

// ConfirmedBooking returns email content or body for the confirmed booking
// of the customer.
func ConfirmedBooking(user schema.User, route schema.BusRoute, booking schema.Bookings, customerSupport string) string {
	var msg string

	msg = fmt.Sprintf("Hello %s,\n", user.FirstName)
	msg += fmt.Sprintf("We are pleased to inform you that your booking from <b>%s</b> to <b>%s</b> on <b>%s</b> has been successfully confirmed.", route.FromRoute, route.ToRoute, booking.TravelDate)
	msg += "&nbsp;Please find below the details of your booking:\n\n"

	msg += bookingDetails(user, route, booking)

	msg += fmt.Sprintf("If you have any questions or clarifications regarding your booking, please feel free to reach out to our customer support team at %s. Thank you and have a pleasant trip!", customerSupport)

	msg = strings.ReplaceAll(msg, "\n", "<br/>")
	msg = strings.ReplaceAll(msg, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;")

	return msg
}
