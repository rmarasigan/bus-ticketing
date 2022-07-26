package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/rmarasigan/bus-ticketing/pkg/common"
	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
)

type Bus struct {
	ID           string `json:"id"`
	Owner        string `json:"owner"`
	Email        string `json:"email,omitempty"`
	Address      string `json:"address"`
	Company      string `json:"company"`
	MobileNumber string `json:"mobile_number"`
	DateCreated  string `json:"date_created"`
}

// SetValues automatically generates the Bus ID as your primary key,
// and set the date it was created.
func (bus *Bus) SetValues() {
	key, err := common.RemoveVowel(bus.Company)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "SetValues", Message: "Failed to remove vowel letters"})
		return
	}

	bus.DateCreated = fmt.Sprint(time.Now().Unix())

	key = strings.Replace(strings.ToUpper(key), " ", "", -1)
	bus.ID = fmt.Sprintf("%s-%s", key, bus.DateCreated[2:8])
}

// ValidateUpdate validates if the field that are going to be updated are empty or not
// to set its previous value.
func (bus *Bus) ValidateUpdate(old *Bus) {
	if bus.Owner == "" {
		bus.Owner = old.Owner
	}

	if bus.Address == "" {
		bus.Address = old.Address
	}

	if bus.Email == "" {
		bus.Email = old.Email
	}

	if bus.MobileNumber == "" {
		bus.MobileNumber = old.MobileNumber
	}
}

type BusUnit struct {
	ID          string `json:"id"`
	Bus         Bus    `json:"bus,omitempty"`
	Code        string `json:"unit_code"`
	Active      bool   `json:"unit_active"`
	Capacity    int    `json:"unit_capacity"`
	DateCreated string `json:"unit_date_created"`
}

type BusRoute struct {
	ID            string  `json:"id"`
	BusUnit       BusUnit `json:"bus_unit,omitempty"`
	Name          string  `json:"route_name"`
	Rate          float64 `json:"route_rate"`
	Available     bool    `json:"route_available"`
	DepartureTime string  `json:"route_departure_time"`
	ArrivalTime   string  `json:"route_arrival_time"`
	FromRoute     string  `json:"route_from"`
	ToRouteCode   string  `json:"route_to"`
	DateCreated   string  `json:"route_date_created"`
}

type BusTrip struct {
	ID          string   `json:"id"`
	BusRoute    BusRoute `json:"bus_route"`
	Status      string   `json:"trip_status"`
	SeatsLeft   int      `json:"trip_seats_left"`
	DateCreated string   `json:"trip_date_created"`
}
