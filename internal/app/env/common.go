package env

import "os"

var (
	BUS_TABLE       = os.Getenv("BUS_TABLE")
	USERS_TABLE     = os.Getenv("USERS_TABLE")
	BUS_UNIT        = os.Getenv("BUS_UNIT_TABLE")
	BUS_ROUTE_TABLE = os.Getenv("BUS_ROUTE_TABLE")
)
