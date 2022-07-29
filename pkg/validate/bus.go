package validate

import (
	"fmt"
	"strings"

	"github.com/rmarasigan/bus-ticketing/pkg/models"
)

// CreateBus validates if the required request parameters are empty or not.
// If some of the fields are empty, it will return an error message.
//
// Required request parameter: company, owner, email, address, mobile_number
func CreateBus(bus *models.Bus) string {
	var msg []string
	var err_msg string

	if bus.Company == "" {
		msg = append(msg, "company")
	}

	if bus.Owner == "" {
		msg = append(msg, "owner")
	}

	if bus.Email == "" {
		msg = append(msg, "email")
	}

	if bus.Address == "" {
		msg = append(msg, "address")
	}

	if bus.MobileNumber == "" {
		msg = append(msg, "mobile_number")
	}

	if len(msg) > 0 {
		err_msg = fmt.Sprintf("Missing %s field(s)", strings.Join(msg, ", "))
	}

	return err_msg
}

// CreateBusUnit checks if the required request parameters are empty or not.
// If some of the fields are empty, it will return an error message.
//
// Required request parameter: code, active, capacity
func CreateBusUnit(unit *models.BusUnit) string {
	var msg []string
	var err_msg string

	if unit.Code == "" {
		msg = append(msg, "code")
	}

	if unit.Active == nil {
		msg = append(msg, "active")
	}

	if unit.Capacity == 0 {
		msg = append(msg, "capacity")
	}

	if len(msg) > 0 {
		err_msg = fmt.Sprintf("Missing %s field(s)", strings.Join(msg, ", "))
	}

	return err_msg
}

// CreateBusRoute checks if the required request parameters are empty or not.
// If some of the fields are emtpy, it will return an error message.
//
// Required request parameter: rate, currency_code, departure_time, arrival_time,
// from_route, to_route, available
func CreateBusRoute(route *models.BusRoute) string {
	var msg []string
	var err_msg string

	if route.Rate <= 0 {
		msg = append(msg, "rate")
	}

	if route.Currency == "" {
		msg = append(msg, "currency_code")
	}

	if route.Available == nil {
		msg = append(msg, "available")
	}

	if route.DepartureTime == "" {
		msg = append(msg, "departure_time")
	}

	if route.ArrivalTime == "" {
		msg = append(msg, "arrival_time")
	}

	if route.FromRoute == "" {
		msg = append(msg, "from_route")
	}

	if route.ToRoute == "" {
		msg = append(msg, "to_route")
	}

	if len(msg) > 0 {
		err_msg = fmt.Sprintf("Missing %s field(s)", strings.Join(msg, ", "))
	}

	return err_msg
}
