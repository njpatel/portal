package shared

import (
	"fmt"
	"os"
)

// ExitOnError - if error is valid, then print the msg and then exit
func ExitOnError(err error, msg string, args ...interface{}) {
	if err != nil {
		fmt.Printf("Error: "+msg+"\n", args...)
		os.Exit(1)
	}
}
