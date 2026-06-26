package main

// cli.go - functions and structs dealing with the command line tool for interacting with tremd
// TODO: Finish the daemon then revisit this mess

import (
	"errors"
	_ "fmt"
	"os"
)

type CommandIdx int
const (
	CIUNSPEC CommandIdx = iota
	CIUNKNOWN
	CIADD
	CISET
	CILIST
	CITOOBIG
)

func CommandIdxToString(input CommandIdx) (string) {
	strings := map[CommandIdx]string{
		CIUNSPEC: "CommandIdx ERROR: Unspecified",
		CIUNKNOWN: "CommandIdx ERROR: Unknown",
		CIADD: "CommandIdx: Add",
		CISET: "CommandIdx: Set",
		CILIST: "CommandIdx: List",
		CITOOBIG: "CommandIdx ERROR: Index OOB",
	}

	return strings[input]
}

type CLIRes struct {
	args []string
	command CommandIdx
	err error
}

func isValidWord(in string, args []string) (func([]string) CLIRes, bool) {
	var cres CLIRes = CLIRes{args: args}
	var wordMap map[string]func([]string)CLIRes = map[string]func([]string)CLIRes{
		"add": 	func([]string) CLIRes {print("got add\n"); cres.command = CIADD; return cres},
		"set": 	func([]string) CLIRes {print("got set\n"); cres.command = CISET; return cres},
		"list": func([]string) CLIRes {print("got list\n"); cres.command = CILIST; return cres},
	}

	res, ok := wordMap[in]
	return res, ok
}

func ProcCLIArgs() CLIRes {
	if len(os.Args) < 2 {return CLIRes{err: errors.New("Too few CLI arguments")}}
	workingArgs := os.Args[1:]
	fp, ok := isValidWord(workingArgs[0], workingArgs)
	if !ok {return CLIRes{err: errors.New(workingArgs[0] + " is not a valid operation")}}

	return fp(workingArgs)
}