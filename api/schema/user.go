package schema

import (
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
	ID           string `json:"id" dynamodbav:"id"`
	UserType     string `json:"user_type" dynamodbav:"user_type"`
	FirstName    string `json:"first_name" dynamodbav:"first_name"`
	LastName     string `json:"last_name" dynamodbav:"last_name"`
	Username     string `json:"username" dynamodbav:"username"`
	Password     string `json:"password,omitempty" dynamodbav:"password"`
	Address      string `json:"address" dynamodbav:"address"`
	Email        string `json:"email" dynamodbav:"email"`
	MobileNumber string `json:"mobile_number" dynamodbav:"mobile_number"`
	DateCreated  string `json:"date_created,omitempty" dynamodbav:"date_created,omitemptyelem"`
	LastLogin    string `json:"last_login,omitempty" dynamodbav:"last_login,omitemptyelem"`
}

// Error sets the default key-value pair.
func (user User) Error(err error, code, message string, kv ...utility.KVP) {
	if user != (User{}) {
		kv = append(kv, utility.KVP{Key: "user", Value: user})
	}

	kv = append(kv, utility.KVP{Key: "Integration", Value: "Bus Ticketing – User"})
	utility.Error(err, code, message, kv...)
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
