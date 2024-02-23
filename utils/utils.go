package utils

import "regexp"

func IsAlpha(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z_]+$`).MatchString(s)
}

func IsNumeric(s string) bool {
	return regexp.MustCompile(`^[0-9]+$`).MatchString(s)
}

func IsSkipable(s string) bool {
	return s == " " || s == "\n" || s == "\t" || s == "\r"
}
