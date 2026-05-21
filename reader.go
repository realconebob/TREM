package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Trem_ReadFileResult struct {
	keyword string
	inputs []string
	callback func([]string) Trem_ReadFileResult
	emsg error
	result string
}

func Trem_ReadFileResult_new(keyword string, inputs []string, callback func([]string) Trem_ReadFileResult, emsg error, result string) *Trem_ReadFileResult {
	if callback == nil {callback = func([]string)Trem_ReadFileResult {return *Trem_ReadFileResult_new("", nil, nil, nil, "")}}

	var res *Trem_ReadFileResult = new(Trem_ReadFileResult)
	res.keyword = keyword
	res.inputs = inputs
	res.callback = callback
	res.emsg = emsg
	res.result = result

	return res
}

func (a *Trem_ReadFileResult) Compare(b Trem_ReadFileResult) bool {
	if a.keyword != b.keyword {return false}
	if a.result != b.result {return false}
	if a.emsg != b.emsg {return false}
	return true
}

func (a *Trem_ReadFileResult) Deep_compare(b Trem_ReadFileResult) bool {
	ares := a.callback(a.inputs)
	bres := b.callback(b.inputs)

	if !ares.Compare(bres) {return false}
	if len(ares.inputs) != len(bres.inputs) {return false}
	for i := 0; i < len(ares.inputs); i++ {
		if ares.inputs[i] != bres.inputs[i] {return false}
	}

	return true
}

func Trem_ReadFile(file *os.File, funcmap map[string]func([]string) Trem_ReadFileResult) ([]Trem_ReadFileResult, error) {
	if file == nil {return nil, errors.New("File is null")}
	if funcmap == nil {return nil, errors.New("Funcmap is null")}

	var res []Trem_ReadFileResult = make([]Trem_ReadFileResult, 0)
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