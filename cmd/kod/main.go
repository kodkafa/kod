package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"kodkafa/internal/app"
	"kodkafa/internal/build"
)

func main() {
	// runtime version injection
	if build.Version == "dev" {
		if out, err := exec.Command("git", "describe", "--tags", "--always", "--dirty").Output(); err == nil {
			build.Version = strings.TrimSpace(string(out))
		}
		if out, err := exec.Command("git", "rev-parse", "HEAD").Output(); err == nil {
			build.Commit = strings.TrimSpace(string(out))
		}
	}

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "KODKAFA Error: %v\n", err)
		os.Exit(1)
	}
}
