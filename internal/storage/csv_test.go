package storage

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func setup(t *testing.T) (*CSVStorage, func()) {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}

	storage := NewCSVStorage(tmpfile.Name())

	return storage, func() {
		os.Remove(tmpfile.Name())
	}
}

func TestCSVStorageAppend(t *testing.T) {
	storage, teardown := setup(t)
	defer teardown()

	data := map[string]string{"email": "example@test.com"}

	t.Run("Append data to storage", func(t *testing.T) {
		if err := storage.Append(data); err != nil {
			t.Fatalf("failed to append data: %v", err)
		}
	})
}

func TestCSVStorageAllRecords(t *testing.T) {
	storage, teardown := setup(t)
	defer teardown()

	data := map[string]string{"email": "example@test.com"}
	if err := storage.Append(data); err != nil {
		t.Fatalf("failed to append data: %v", err)
	}

	t.Run("Read data from storage", func(t *testing.T) {
		readData, err := storage.AllRecords()
		if err != nil {
			t.Fatalf("failed to read data: %v", err)
		}

		if diff := cmp.Diff(data, readData[0]); diff != "" {
			t.Errorf("read data does not match written data (-want +got):\n%s", diff)
		}
	})
}
