package reminders

import (
	"github.com/realconebob/trem/internal"
	"encoding/gob"
	"bytes"
	"os"
)

func SerializeToGobFile[T any](entries []T, path string, filter func(T)bool) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	working := misc.Filter(entries, filter)
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(working); err != nil {
		return err
	}
	return nil
}

func GetFromGob[T any](data []byte) ([]T, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	var reminders []T

	if err := decoder.Decode(&reminders); err != nil {
		return []T{}, err
	}

	return reminders, nil
}

func GetFromGobFile[T any](path string) ([]T, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return []T{}, err
	}
	return GetFromGob[T](contents)
}
