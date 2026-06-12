package main

import (
	"fmt"
	"os"
)

func main() {
	res := ProcCLIArgs()
	if res.err != nil {fmt.Fprintf(os.Stderr, "Encountered an error: %v\n", res.err); os.Exit(1)}

	print("args: ")
	for _, arg := range res.args {
		print(arg + " ")
	}
	print("\n")
}