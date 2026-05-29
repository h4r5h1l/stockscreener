package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("Running Python from Go...!")
	cmd := exec.Command("uv", "run", "../python/main.py")

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(string(out))
}   