package common

import (
	"regexp"

	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
)

// RemoveVowel removes vowel letters and some characters (., *, (, ), \)
// from the string passed and returns a string without the vowel letters.
//
// Example:
//  input, err := RemoveVowel(`a(bc\de.f)g*hi`)
//  if err != nil {
// 	log.Print(err)
//  }
//  fmt.Println(input)
//
// Output:
//  bcdfgh
func RemoveVowel(str string) (string, error) {
	regex, err := regexp.Compile(`[aeiouyAEIOUY\.\*\_\(\)\\]`)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "RemoveVowel", Message: "Failed to parse regex expression"})

		return "", err
	}

	return regex.ReplaceAllString(str, ""), nil
}
