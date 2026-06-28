package main

import (
	"fmt"
	"os"
)

func main() {
	res := ProcCLIArgs()
	if res.err != nil {fmt.Fprintf(os.Stderr, "Encountered an error: %v\n", res.err); os.Exit(1)}

	switch res.command {
	case TC_ADD: 	print("call adding behavior here\n")
	case TC_EDIT: 	print("call setting behavior here\n")
	case TC_LIST: 	print("call listing behavior here\n")
	case TC_REMOVE: print("call removing behavior here\n")
	default: 		fmt.Fprintf(os.Stderr, "Encountered an invalid command: %v\n", res.command)
	}
}