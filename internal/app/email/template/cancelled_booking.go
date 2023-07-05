package template

import (
	"fmt"
	"strings"

	"github.com/rmarasigan/bus-ticketing/api/schema"
)

// CustomerCancelledBooking returns e-mail content for the canceled booking
// made by the customer.
func CustomerCancelledBooking(user schema.User, route schema.BusRoute, booking schema.Bookings, customerSupport string) string {
	var msg string

	msg = fmt.Sprintf("Hello %s,\n", user.FirstName)
	msg += "We have received your request to cancel your booking with the following details:\n\n"

	msg += bookingDetails(user, route, booking)

	msg += "We have processed your cancellation request, and we confirm that your booking has been successfully canceled as per your instructions.\n"
	msg += fmt.Sprintf("If you have any further questions or require assistance, please feel free to contact our customer support team at %s.\n", customerSupport)

	msg = strings.ReplaceAll(msg, "\n", "<br/>")
	msg = strings.ReplaceAll(msg, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;")

	return msg
}

// CancelledBooking returns e-mail content for the canceled booking
// made by the administrator.
func CancelledBooking(user schema.User, route schema.BusRoute, booking schema.Bookings, customerSupport string) string {
	var msg string

	msg = fmt.Sprintf("Hello %s,\n", user.FirstName)
	msg += "We regret to inform you that, due to unforeseen circumstances beyond our control, we must cancel your bus booking with the following details:\n\n"

	msg += bookingDetails(user, route, booking)

	msg += "We apologize for any inconvenience caused by this cancellation, and we understand the impact it may have on your travel plans. Rest assured, our team is working diligently to address the situation and explore alternative solutions.\n\n"
	msg += fmt.Sprintf("If you have any further questions or require assistance, please feel free to contact our customer support team at %s.\n", customerSupport)

	msg = strings.ReplaceAll(msg, "\n", "<br/>")
	msg = strings.ReplaceAll(msg, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;")

	return msg
}
