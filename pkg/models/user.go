package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
)

const (
	ADMN  = "ADMIN"
	CSTMR = "CUSTOMER"
)

var UserType = map[int]string{
	1: ADMN,
	2: CSTMR,
}

var UserIDCode = map[int]string{
	1: "ADMN",
	2: "CSTMR",
}

type User struct {
	ID           string `json:"id"`
	UserType     string `json:"type"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile_number,omitempty"`
	DateCreated  string `json:"date_created"`
}

type UserResponse struct {
	ID           string `json:"id"`
	UserType     string `json:"type"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile_number,omitempty"`
}

// SetValues automatically generates the User ID as your primary key,
// set the user type and the date it was created.
func (user *User) SetValues() {
	Type, err := strconv.Atoi(user.UserType)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "SetValues", Message: "Failed to convert type string to int"})
		return
	}

	user.UserType = UserType[Type]
	user.DateCreated = fmt.Sprint(time.Now().Unix())
	user.ID = fmt.Sprintf("%s-%s", UserIDCode[Type], user.DateCreated[2:8])
}

// ValidateUpdateAccount validates if the field that are going to be updated are empty or not
// to set its previous value.
func (user *UserResponse) ValidateUpdateAccount(old *UserResponse) *UserResponse {
	if user.FirstName == "" {
		user.FirstName = old.FirstName
	}

	if user.LastName == "" {
		user.LastName = old.LastName
	}

	if user.Address == "" {
		user.Address = old.Address
	}

	if user.Email == "" {
		user.Email = old.Email
	}

	if user.MobileNumber == "" {
		user.MobileNumber = old.MobileNumber
	}

	if strings.HasPrefix(user.ID, "ADMN") {
		user.UserType = "ADMIN"
	}

	if strings.HasPrefix(user.ID, "CSTMR") {
		user.UserType = "CUSTOMER"
	}

	return user
}
