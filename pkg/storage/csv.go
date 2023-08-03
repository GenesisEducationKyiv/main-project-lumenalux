package storage

import (
	"encoding/csv"
	"os"
)

var _headers = []string{"email"} // The order of the columns keys

type StorageConfig struct {
	Path string `default:"./storage/storage.csv"`
}

type CSVStorage struct {
	FilePath string
}

func NewCSVStorage(filePath string) *CSVStorage {
	return &CSVStorage{FilePath: filePath}
}

func (s *CSVStorage) AllRecords() ([]map[string]string, error) {
	f, err := os.Open(s.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var maps []map[string]string

	for _, record := range records {
		rowMap := make(map[string]string, len(_headers))
		for i, key := range _headers {
			rowMap[key] = record[i]
		}
		maps = append(maps, rowMap)
	}

	return maps, nil
}

func (s *CSVStorage) Append(record map[string]string) error {
	f, err := os.OpenFile(s.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)

	// Build a slice of values based on the order of the keys
	var values []string
	for _, key := range _headers {
		values = append(values, record[key])
	}

	if err = w.Write(values); err != nil {
		return err
	}
	w.Flush()

	return w.Error()
}
