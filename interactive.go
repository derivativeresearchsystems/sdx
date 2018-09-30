package main

import (
	"bufio"
	"fmt"
	"os"
)

// manualoverride allows for using keystrokes to override the
// current context
func manualoverride(clocal, cremote string, ccurent *string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "l", "local", "use local":
			displayinfo(fmt.Sprintf("Overriding state, switching to local context %v", clocal))
			ccurent = &clocal
		case "r", "remote", "use remote":
			displayinfo(fmt.Sprintf("Overriding state, switching to remote context %v", cremote))
			ccurent = &cremote
		default:
		}
	}
}
