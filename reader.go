package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// Classes and functions to read and interpret files

type AOFResult[T any] struct {
	line string
	oper func(string) T
	res T
}

func (s *AOFResult[T]) _set_empty() {
	s.line = ""
	s.oper = func(_ string) T {return *new(T)}
	s.res = *new(T)
}

func New_AOFResult[T any]() *AOFResult[T] {
	res := new(AOFResult[T])
	res._set_empty()
	return res
}

func ActOnFile[T any](path string, delim string, funcmap map[string]func(string) T) ([]AOFResult[T], error) {
	if len(path) < 1 	{return nil, errors.New("path is undefined")}
	if len(funcmap) < 1 {return nil, errors.New("funcmap is empty/nil")}

	res := make([]AOFResult[T], 0);
	filedata, err := os.Open(path)
	if err != nil {return nil, errors.Join(errors.New("could not read file \"" + path + "\""), err)}
	defer filedata.Close()
	scanner := bufio.NewScanner(filedata)

	// Read file line by line, checking if the line starts with a keyword. If so, process it
	for scanner.Scan() {
		tmp := *New_AOFResult[T]()

		tmp.line = scanner.Text()
		lines := strings.Split(tmp.line, delim)
		if len(lines) != 2 {
			res = append(res, tmp)
			continue
		}

		oper, exists := funcmap[lines[0]]
		if !exists {
			res = append(res, tmp)
			continue
		}
		tmp.oper = oper

		tmp.res = tmp.oper(tmp.line)
		res = append(res, tmp)
	}

	return res, nil;
}