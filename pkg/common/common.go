package common

import (
	"regexp"

	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
)

// RemoveVowel removes vowel letters from the string passed
// and returns a string without the vowel letters.
func RemoveVowel(str string) (string, error) {
	regex, err := regexp.Compile(`[aeiouyAEIOUY]`)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "RemoveVowel", Message: "Failed to parse regex expression"})

		return "", err
	}

	return regex.ReplaceAllString(str, ""), nil
}

// RemoveSymbols removes symbols from the string passed and
// returns a string without the symbols and whitespaces.
func RemoveSymbols(str string) (string, error) {
	regex, err := regexp.Compile(`[^\w]`)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "RemoveSymbols", Message: "Failed to parse regex expression."})
		return "", err
	}

	return regex.ReplaceAllString(str, ""), nil

}
