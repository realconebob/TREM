package cli

// cli.go - functions and structs dealing with the command line tool for interacting with tremd
// TODO: Finish the daemon then revisit this mess

import (
	"errors"
	"fmt"
	"os"
)

type TremCommand uint32

const (
	TC_UNSPEC TremCommand = iota
	TC_UNKNOWN

	TC_LIST
	TC_ADD
	TC_EDIT
	TC_REMOVE
	TC_DAEMON

	TC_TOOBIG
)

type CLIRes struct {
	Command   TremCommand
	Arguments []string
	Err       error
}

func CreateCLIResFromArgs(command string, args []string) CLIRes {
	var res CLIRes = CLIRes{Arguments: args}
	switch command {
	case "add":
		res.Command = TC_ADD
	case "edit":
		res.Command = TC_EDIT
	case "list":
		res.Command = TC_LIST
	case "remove":
		res.Command = TC_REMOVE
	case "daemon":
		res.Command = TC_DAEMON
	default:
		res.Err = errors.New("Unknown command \"" + command + "\"")
		res.Command = TC_UNKNOWN
	}

	return res
}

func ProcCLIArgs() CLIRes {
	fmt.Println("Arguments:")
	for _, arg := range os.Args {
		fmt.Printf("\"%v\" ", arg)
	}
	fmt.Print("\n")

	if len(os.Args) < 2 {
		return CLIRes{Err: errors.New("Too few CLI arguments")}
	}
	workingArgs := os.Args[1:]
	return CreateCLIResFromArgs(workingArgs[0], workingArgs[1:])
}
