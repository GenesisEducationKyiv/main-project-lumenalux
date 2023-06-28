package storage

import (
	"os"
	"reflect"
	"testing"
)

func TestCSVStorage(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}

	defer os.Remove(tmpfile.Name())

	storage := NewCSVStorage(tmpfile.Name())

	data := "example@test.com"

	if err := storage.Append(data); err != nil {
		t.Fatalf("failed to append data: %v", err)
	}

	readData, err := storage.AllRecords()
	if err != nil {
		t.Fatalf("failed to read data: %v", err)
	}

	if !reflect.DeepEqual(readData[0][0], data) {
		t.Errorf("read data does not match written data: got %v, want %v", readData[0][0], data)
	}
}
