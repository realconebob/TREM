package misc

import (
	"encoding/gob"
	"bytes"
	"os"
	"io"
)

func serializeToWriter[T any](writer io.Writer, entries []T, filter func(T)bool) error {
	encoder := gob.NewEncoder(writer)
	working := Filter(entries, filter)

	if err := encoder.Encode(working); err != nil {return err}
	return nil
}

func SerializeToGobFile[T any](entries []T, path string, filter func(T)bool) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return serializeToWriter(file, entries, filter)
}

func SerializeToBuffer[T any](entries []T, filter func(T)bool) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := serializeToWriter(buf, entries, filter)
	if err != nil {return nil, err}

	return buf.Bytes(), nil
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
