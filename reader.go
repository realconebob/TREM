package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type ReadFileResult struct {
	keyword string
	inputs []string
	callback func([]string) ReadFileResult
	emsg error
	result string
}

func ReadFileResult_new(keyword string, inputs []string, callback func([]string) ReadFileResult, emsg error, result string) *ReadFileResult {
	if callback == nil {callback = func([]string)ReadFileResult {return *ReadFileResult_new("", nil, nil, nil, "")}}

	var res *ReadFileResult = new(ReadFileResult)
	res.keyword = keyword
	res.inputs = inputs
	res.callback = callback
	res.emsg = emsg
	res.result = result

	return res
}

func (a *ReadFileResult) Compare(b ReadFileResult) bool {
	if a.keyword != b.keyword {return false}
	if a.result != b.result {return false}
	if a.emsg != b.emsg {return false}
	return true
}

func (a *ReadFileResult) Deep_compare(b ReadFileResult) bool {
	ares := a.callback(a.inputs)
	bres := b.callback(b.inputs)

	if !ares.Compare(bres) {return false}
	if len(ares.inputs) != len(bres.inputs) {return false}
	for i := 0; i < len(ares.inputs); i++ {
		if ares.inputs[i] != bres.inputs[i] {return false}
	}

	return true
}

func ReadFile(file *os.File, funcmap map[string]func([]string) ReadFileResult) ([]ReadFileResult, error) {
	if file == nil {return nil, errors.New("File is null")}
	if funcmap == nil {return nil, errors.New("Funcmap is null")}

	var res []ReadFileResult = make([]ReadFileResult, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		tokens := strings.Split(text, "=")
		if len(tokens) < 2 {continue} // No =, don't care about line contents

		cb, ok := funcmap[tokens[0]]
		if ok && cb != nil {res = append(res, cb(tokens))}
	}

	return res, nil
}