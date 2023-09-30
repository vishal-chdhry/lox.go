package errors

import "fmt"

func Error(line int, where, message string) error {
	msg := fmt.Sprintf("[line %d] Error %s: %s", line, where, message)
	report(msg)
	return fmt.Errorf(msg)
}

func report(message string) {
	fmt.Println(message)
}
