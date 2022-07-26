package validate

import (
	"fmt"
	"strings"

	"github.com/rmarasigan/bus-ticketing/pkg/models"
)

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
