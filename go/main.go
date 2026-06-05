package main

import (
	"fmt"
	"github.com/h4r5h1l/stockscreener/handler" // Import your cmd package
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

	fmt.Println("Running IngestEquities from Go...!")
	handler.IngestEquities()

}
