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
		msg = append(msg, "Company")
	}

	if bus.Owner == "" {
		msg = append(msg, "Owner")
	}

	if bus.Email == "" {
		msg = append(msg, "Email")
	}

	if bus.Address == "" {
		msg = append(msg, "Address")
	}

	if bus.MobileNumber == "" {
		msg = append(msg, "MobileNumber")
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
		msg = append(msg, "Code")
	}

	if unit.Active == nil {
		msg = append(msg, "Active")
	}

	if unit.Capacity == 0 {
		msg = append(msg, "Capacity")
	}

	if len(msg) > 0 {
		err_msg = fmt.Sprintf("Missing %s field(s)", strings.Join(msg, ", "))
	}

	return err_msg
}

// CreateBusRoute checks if the required request parameters are empty or not.
// If some of the fields are emtpy, it will return an error message.
//
// Required request parameter: rate, departure_time, arrival_time, from_route
// to_route, available
func CreateBusRoute(route *models.BusRoute) string {
	var msg []string
	var err_msg string

	if route.Rate <= 0 {
		msg = append(msg, "Rate")
	}

	if route.Available == nil {
		msg = append(msg, "Available")
	}

	if route.DepartureTime == "" {
		msg = append(msg, "DepartureTime")
	}

	if route.ArrivalTime == "" {
		msg = append(msg, "ArrivalTime")
	}

	if route.FromRoute == "" {
		msg = append(msg, "FromRoute")
	}

	if route.ToRoute == "" {
		msg = append(msg, "ToRoute")
	}

	if len(msg) > 0 {
		err_msg = fmt.Sprintf("Missing %s field(s)", strings.Join(msg, ", "))
	}

	return err_msg
}
