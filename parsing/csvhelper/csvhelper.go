package csvhelper

import (
	"encoding/csv"
	"io"
	"os"
)

func ParseCSVFile(filepath string) ([][]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	parsed, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func WriteCSVFile(filepath string, data [][]string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	return WriteCSV(f, data)
}

func WriteCSV(w io.Writer, data [][]string) error {
	cw := csv.NewWriter(w)
	return cw.WriteAll(data)
}
