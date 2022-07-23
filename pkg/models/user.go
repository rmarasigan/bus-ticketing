package models

import (
	"fmt"
	"strconv"
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
