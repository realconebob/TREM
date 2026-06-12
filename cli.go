package main

import (
	"errors"
	_ "fmt"
	"os"
)

type CLIRes struct {
	args []string
	err error
}

func isValidWord(in string, args []string) (func([]string) CLIRes, bool) {
	var cres CLIRes = CLIRes{args: args}
	var wordMap map[string]func([]string)CLIRes = map[string]func([]string)CLIRes{
		"add": 	func([]string) CLIRes {print("got add\n"); return cres},
		"set": 	func([]string) CLIRes {print("got set\n"); return cres},
		"list": func([]string) CLIRes {print("got list\n"); return cres},
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