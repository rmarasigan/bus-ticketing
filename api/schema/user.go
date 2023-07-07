package schema

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/rmarasigan/bus-ticketing/internal/utility"
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

// User contains the user account information.
//
// The "dynamodbav" struct tag can be used to control the value
// that will be marshaled into a AttributeValue.
type User struct {
	ID           string `json:"id" dynamodbav:"id"`                                             // The unique user ID and the sort key
	UserType     string `json:"user_type" dynamodbav:"user_type"`                               // The type of the user account (either ADMIN or CUSTOMER)
	FirstName    string `json:"first_name" dynamodbav:"first_name"`                             // The first name of the user
	LastName     string `json:"last_name" dynamodbav:"last_name"`                               // The last name of the user
	Username     string `json:"username" dynamodbav:"username"`                                 // The username of the user account and the primary key
	Password     string `json:"password,omitempty" dynamodbav:"password"`                       // THe user security password for the account
	Address      string `json:"address" dynamodbav:"address"`                                   // The user address
	Email        string `json:"email" dynamodbav:"email"`                                       // The user e-mail address
	MobileNumber string `json:"mobile_number" dynamodbav:"mobile_number"`                       // The user phone
	DateCreated  string `json:"date_created,omitempty" dynamodbav:"date_created,omitemptyelem"` // The date it was created
	LastLogin    string `json:"last_login,omitempty" dynamodbav:"last_login,omitemptyelem"`     // The last login session of the user
}

// Error sets the default key-value pair.
func (user User) Error(err error, code, message string, kv ...utility.KVP) {
	if user != (User{}) {
		kv = append(kv, utility.KVP{Key: "user", Value: user})
	}

	kv = append(kv, utility.KVP{Key: "Integration", Value: "Bus Ticketing – User"})
	utility.Error(err, code, message, kv...)
}

// IsEmptyPayload checks if the request payload is empty and if it is,
// it will return an error message.
func (user User) IsEmptyPayload(payload string) error {
	if payload == "" {
		err := errors.New("payload is required")
		user.Error(err, "APIError", "the request payload is empty")

		return err
	}

	return nil
}

// LastLogIn returns the current time when the user logged in.
func (user User) LastLogIn() string {
	return time.Now().Format("02 Jan 2006 15:04:05")
}

// SetValues automatically generates the User ID as your primary key,
// set the user type and the date it was created.
//
// Example:
//		id: CSTMR-854980
//		user_type: CUSTOMER
//		date_created: 1685498070
func (user *User) SetValues() {
	Type, err := strconv.Atoi(user.UserType)
	if err != nil {
		user.Error(err, "SetValues", "failed to convert type string to int")
		return
	}

	user.UserType = UserType[Type]
	user.DateCreated = fmt.Sprint(time.Now().Unix())
	user.ID = fmt.Sprintf("%s-%s", UserIDCode[Type], user.DateCreated[2:8])
}
