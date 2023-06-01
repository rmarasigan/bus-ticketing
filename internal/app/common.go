package app

import "regexp"

// RemoveVowel removes vowel letters from the string passed
// and returns a string without the vowel letters.
func RemoveVowel(str string) (string, error) {
	regex, err := regexp.Compile(`[aeiouyAEIOUY]`)
	if err != nil {
		return "", err
	}

	return regex.ReplaceAllString(str, ""), nil
}

// RemoveSymbols removes symbols from the string passed and
// returns a string without the symbols and whitespaces.
func RemoveSymbols(str string) (string, error) {
	regex, err := regexp.Compile(`[^\w]`)
	if err != nil {
		return "", err
	}

	return regex.ReplaceAllString(str, ""), nil
}

// IsNumberOnly checks whether the input string is a number. It accepts
// positive fractional number, including zero (0, 1, 0.0, 1.0, 99999.000001,
// 5.1).
func IsNumberOnly(str string) (bool, error) {
	regex, err := regexp.Compile(`^(0|[1-9]\d*)(\.\d+)?$`)
	if err != nil {
		return false, err
	}

	return regex.MatchString(str), nil
}
