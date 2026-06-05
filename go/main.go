package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("Running Python from Go...!")

	// 1. Point directly to the python binary inside your custom 'twsenv' virtual environment
	// 2. Fix the folder name from '../python/' to '../ibpython/'
// The directory path and the script name must be separate arguments
	cmd := exec.Command("uv", "run", "--directory", "../ibpython", "main.py", "fetch_universe") // Correct
	// Captures both standard output and error messages together
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Runtime Execution Error:", err)
	}

	fmt.Println("--- Subprocess Output ---")
	fmt.Println(string(out))
}
