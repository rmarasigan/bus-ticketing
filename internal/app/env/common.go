package env

import "os"

var (
	BUS_TABLE   = os.Getenv("BUS_TABLE")
	USERS_TABLE = os.Getenv("USERS_TABLE")
)
