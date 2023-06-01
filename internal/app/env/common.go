package env

import "os"

var (
	USERS_TABLE = os.Getenv("USERS_TABLE")
)
