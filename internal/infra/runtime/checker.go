package runtime

import (
	"os/exec"
)

// CheckInterpreter checks if the given interpreter exists in the system.
// Returns the path if found, or an empty string and false.
func CheckInterpreter(name string) (string, bool) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", false
	}
	return path, true
}
