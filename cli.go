package main

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
	command TremCommand
	arguments []string
	err error
}

func CreateCLIResFromArgs(command string, args []string) CLIRes {
	var res CLIRes = CLIRes{arguments: args}
	switch command {
	case "add": 	res.command = TC_ADD
	case "edit": 	res.command = TC_EDIT
	case "list": 	res.command = TC_LIST
	case "remove": 	res.command = TC_REMOVE
	case "daemon": 	res.command = TC_DAEMON
	default:
		res.err = errors.New("Unknown command \"" + command + "\"")
		res.command = TC_UNKNOWN
	}

	return res
}

func ProcCLIArgs() CLIRes {
	fmt.Println("Arguments:")
	for _, arg := range os.Args {
		fmt.Printf("\"%v\" ", arg)
	}
	fmt.Print("\n")

	if len(os.Args) < 2 {return CLIRes{err: errors.New("Too few CLI arguments")}}
	workingArgs := os.Args[1:]
	return CreateCLIResFromArgs(workingArgs[0], workingArgs[1:])
}