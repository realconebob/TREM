package main

// misc.go - Miscellaneous functions and structs that don't deserve their own file

import (
	"errors"
	"os"
	"time"
)

func Filter[T any](slice []T, checker func(T) bool) []T {
	var res []T = make([]T, 0)
	for _, entry := range slice {
		if checker(entry) {
			res = append(res, entry)
		}
	}
	return res
}
// Why is there not a filter function in the stdlib? I really shouldn't have to write this myself,
// even if it's trivial. ALSO why has go not implemented receivers on generics? That way the signature
// could be "func (slice *[]T) Filter[T any](checker func(T) bool) *[]T"  could be written like
// "(&[]int{1, 2, 3}).Filter(checker)" or "&sliceToFilter.Filter(checker)" instead of how it is now

type WatchedFile struct {
	Handle *os.File
	LastStat os.FileInfo
	PollingInterval time.Duration
	closed bool
}

func GetFileWatch(path string, polling time.Duration) (*WatchedFile, error) {
	file, err := os.Open(path)
	if err != nil {return &WatchedFile{}, err}

	stat, err := file.Stat()
	if err != nil {return &WatchedFile{}, err}

	return &WatchedFile{
		Handle: file,
		LastStat: stat,
		PollingInterval: polling,
		closed: false,
	}, nil
}

func (watch *WatchedFile) Close() {
	if watch.closed {return}
	watch.Handle.Close()
	watch.closed = true
}

func (watch *WatchedFile) CheckForUpdate() (bool, error) {
	if watch.closed {return false, errors.New("underlying file has been closed")}
	var res bool = false

	newStat, err := watch.Handle.Stat()
	if err != nil {return false, err}

	if newStat.Size() != watch.LastStat.Size() || newStat.ModTime() != watch.LastStat.ModTime() {
		res = true
		watch.LastStat = newStat
	}
	return res, nil
}