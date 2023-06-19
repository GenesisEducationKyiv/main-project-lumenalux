package storage

import (
	"encoding/csv"
	"os"
)

type CSVStorage struct {
	FilePath string
}

func NewCSVStorage(filePath string) *CSVStorage {
	return &CSVStorage{FilePath: filePath}
}

func (s *CSVStorage) Read() ([][]string, error) {
	f, err := os.Open(s.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	return r.ReadAll()
}

func (s *CSVStorage) Append(record []string) error {
	f, err := os.OpenFile(s.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err = w.Write(record); err != nil {
		return err
	}
	w.Flush()

	return w.Error()
}
