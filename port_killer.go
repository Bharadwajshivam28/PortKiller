package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run port_killer.go <port_number>")
		os.Exit(1)
	}

	portStr := os.Args[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("Invalid port number:", err)
		os.Exit(1)
	}

	err = killProcessesByPort(port)
	if err != nil {
		fmt.Println("Error killing processes:", err)
		os.Exit(1)
	}

	fmt.Printf("Processes using port %d forcefully terminated.\n", port)
}

func killProcessesByPort(port int) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", "netstat -ano | findstr LISTENING | findstr "+strconv.Itoa(port))
	case "darwin", "linux":
		cmd = exec.Command("lsof", "-t", "-i", fmt.Sprintf(":%d", port))
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to find processes: %v, output: %s", err, string(output))
	}

	pids := strings.Fields(string(output))
	for _, pid := range pids {
		killCmd := exec.Command("kill", "-9", pid)
		killOutput, killErr := killCmd.CombinedOutput()
		if killErr != nil {
			return fmt.Errorf("failed to kill process %s: %v, output: %s", pid, killErr, string(killOutput))
		}
	}

	return nil
}
