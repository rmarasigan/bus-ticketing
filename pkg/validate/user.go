package validate

import (
	"fmt"
	"strings"

	"github.com/rmarasigan/bus-ticketing/pkg/models"
)

// CreateAccount validates if the fields are empty or not.
func CreateAccount(user *models.User) string {
	var msg []string
	var err_msg string

	if user.UserType == "" {
		msg = append(msg, "UserType")
	}

	if user.FirstName == "" {
		msg = append(msg, "FirstName")
	}

	if user.LastName == "" {
		msg = append(msg, "LastName")
	}

	if user.Username == "" {
		msg = append(msg, "Username")
	}

	if user.Password == "" {
		msg = append(msg, "Password")
	}

	if user.Address == "" {
		msg = append(msg, "Address")
	}

	if user.Email == "" {
		if user.MobileNumber == "" {
			msg = append(msg, "Email, MobileNumber")
		}
	}

	if len(msg) != 0 {
		err_msg = fmt.Sprintf("Missing %s field(s)", strings.Join(msg, ", "))
	}

	return err_msg
}

// LogIn validates if the required fields are empty.
func LogIn(user *models.User) string {
	var msg []string
	var err_msg string

	if user.Username == "" {
		msg = append(msg, "Username")
	}

	if user.Password == "" {
		msg = append(msg, "Password")
	}

	if len(msg) > 0 {
		err_msg = fmt.Sprintf("Missing %s field(s)", strings.Join(msg, ", "))
	}

	return err_msg
}
