package helper

import "strings"

// remove punctiation marks and symbols from string
func RemoveSymbols(str string) string {
	stripped_str := str
	punctuations := "!()-[]{};:'\",<>./?@#$%^&*_~«»"
	for _, char := range str {
		for _, punctuation := range punctuations {
			if char == punctuation {
				stripped_str = strings.Replace(stripped_str, string(char), "", -1)
			}
		}
	}
	return stripped_str
}