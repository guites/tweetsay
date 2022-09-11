package helper

import "fmt"

func Prompt_key_from_user(keyName string) (string) {
	var keyValue string
	fmt.Printf("Please enter you %s: ", keyName)
	fmt.Scanf("%s", &keyValue)
	return keyValue
}