package misc

// misc.go - Miscellaneous functions and structs that don't deserve their own file

import (
	"errors"
	"math"
	"time"
	"fmt"
	"os"
)

// Returns a new array from the provided array, excluding any elements that checker returns false for
func Filter[T any](slice []T, checker func(T) bool) []T {
	var res []T = make([]T, 0)
	for _, entry := range slice {
		if checker(entry) {
			res = append(res, entry)
		}
	}
	return res
}

func PseudoRandom(feed uint64) uint64 {
	return (feed * 10102007) ^ math.Float64bits(math.Exp(math.Float64frombits(feed))) ^ math.Float64bits(math.Pow(2, 0.5)) ^ feed
}

func PrintErrAndExit(err error, format string, a ...any) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, a...)
		os.Exit(1)
	}
}

type Comparable[T any] interface {
	Compare(T) bool
}

func IsListEqual[T Comparable[T]](l1, l2 []T) bool {
	if len(l1) != len(l2) {return false}
	for idx := range l1 {
		if !l1[idx].Compare(l2[idx]) {return false}
	}

	return true
}

// Runs through a slice of errors to check if they're non-nil. If an entry is non-nil, it returns the list. Otherwise, returns nil
func NilErrSliceCheck(eslice []error) []error {
	for _, e := range eslice {
		if e != nil {return eslice}
	}
	return nil
}

type WatchedFile struct {
	Handle          *os.File
	LastStat        os.FileInfo
	PollingInterval time.Duration
	closed          bool
}

func GetFileWatch(path string, polling time.Duration) (*WatchedFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return &WatchedFile{}, err
	}

	stat, err := file.Stat()
	if err != nil {
		return &WatchedFile{}, err
	}

	return &WatchedFile{
		Handle:          file,
		LastStat:        stat,
		PollingInterval: polling,
		closed:          false,
	}, nil
}

func (watch *WatchedFile) Close() error {
	if watch == nil {return errors.New("watch is nil")}
	if watch.closed {
		return nil
	}
	err := watch.Handle.Close()
	watch.closed = true
	return err
}

func (watch *WatchedFile) CheckForUpdate() (bool, error) {
	if watch == nil {return false, errors.New("watch is nil")}
	if watch.closed {
		return false, errors.New("underlying file has been closed")
	}
	var res bool = false

	newStat, err := watch.Handle.Stat()
	if err != nil {
		return false, err
	}

	if newStat.Size() != watch.LastStat.Size() || newStat.ModTime() != watch.LastStat.ModTime() {
		res = true
		watch.LastStat = newStat
	}
	return res, nil
}

func (watch *WatchedFile) Sleep() {
	if watch == nil {return}
	time.Sleep(watch.PollingInterval)
}