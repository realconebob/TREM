package main

import (
	"fmt"
	"os"
)

func main() {
	res := ProcCLIArgs()
	if res.err != nil {fmt.Fprintf(os.Stderr, "Encountered an error: %v\n", res.err); os.Exit(1)}

	switch res.command {
		case CIADD: {print("call adding behavior here\n")}
		case CISET: {print("call setting behavior here\n")}
		case CILIST: {print("call listing behavior here\n")}
		default: {
			fmt.Fprintf(os.Stderr, "Encountered an invalid command: %v\n", CommandIdxToString(res.command))
		}
	}
}