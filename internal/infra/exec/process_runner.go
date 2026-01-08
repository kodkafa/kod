package exec

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"kodkafa/internal/domain/entities"
	"kodkafa/internal/domain/ports"
)

type ProcessRunner struct {
	baseDir string
}

func NewProcessRunner(baseDir string) *ProcessRunner {
	return &ProcessRunner{baseDir: baseDir}
}

func (r *ProcessRunner) Run(plugin *entities.Plugin, args string, mode ports.RunMode, outputChan chan<- ports.OutputChunk) (*ports.RunResult, error) {
	start := time.Now()
	var outputBuilder strings.Builder

	var cmd *exec.Cmd
	cmdArgs := parseArgs(args)

	switch plugin.Interpreter {
	case "python":
		pythonPath := filepath.Join(r.baseDir, "core", "python", "venv", "bin", "python3")
		entryPath := filepath.Join(plugin.Source, plugin.Entry)
		fullArgs := append([]string{entryPath}, cmdArgs...)
		cmd = exec.Command(pythonPath, fullArgs...)
	case "node":
		nodeCoreDir := filepath.Join(r.baseDir, "core", "node")
		entryPath := filepath.Join(plugin.Source, plugin.Entry)
		fullArgs := append([]string{entryPath}, cmdArgs...)
		cmd = exec.Command("node", fullArgs...)
		// Set NODE_PATH to use core/node/node_modules
		cmd.Env = append(os.Environ(), "NODE_PATH="+filepath.Join(nodeCoreDir, "node_modules"))
	case "r":
		entryPath := filepath.Join(plugin.Source, plugin.Entry)
		fullArgs := append([]string{entryPath}, cmdArgs...)
		cmd = exec.Command("Rscript", fullArgs...)
	default:
		return nil, fmt.Errorf("unsupported interpreter: %s", plugin.Interpreter)
	}

	cmd.Dir = plugin.Source

	if outputChan != nil {
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		// Start and streaming logic remains same...
		_ = cmd.Start()

		go func() {
			defer close(outputChan)
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				text := scanner.Text()
				outputChan <- ports.OutputChunk{Data: scanner.Bytes(), Plugin: plugin.Name}
				outputBuilder.WriteString(text + "\n")
			}
			scannerErr := bufio.NewScanner(stderr)
			for scannerErr.Scan() {
				text := scannerErr.Text()
				outputChan <- ports.OutputChunk{Data: scannerErr.Bytes(), IsErr: true, Plugin: plugin.Name}
				outputBuilder.WriteString(text + "\n")
			}
		}()

		err := cmd.Wait()
		exitCode := 0
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			} else {
				exitCode = 1
			}
		}

		return &ports.RunResult{
			ExitCode: exitCode,
			Duration: time.Since(start).Nanoseconds(),
			Status:   "completed",
			Output:   outputBuilder.String(),
		}, nil
	}

	// CLI Direct mode
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return &ports.RunResult{
		ExitCode: exitCode,
		Duration: time.Since(start).Nanoseconds(),
		Status:   "completed",
	}, nil
}

// parseArgs parses a command line string into arguments, respecting quotes.
func parseArgs(args string) []string {
	var parts []string
	var current string
	var inQuote bool
	var quoteChar rune

	for _, r := range args {
		switch {
		case r == ' ' && !inQuote:
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		case (r == '"' || r == '\'') && !inQuote:
			inQuote = true
			quoteChar = r
		case r == quoteChar && inQuote:
			inQuote = false
		default:
			current += string(r)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
