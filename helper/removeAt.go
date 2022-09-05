package helper

// remove @ if user included in command line arg
func RemoveAt(username string) (string) {
	if username[0:1] == "@" {
		username = username[1:]
	}
	return username
}