package commands

import (
	"encoding/csv"
	"os"
	"strconv"
)

func readCsv(path string) (list []SendMessage, err error) {
	file, err := os.Open(path)
	if err != nil {
		return []SendMessage{}, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		if len(record) != 2 {
			return []SendMessage{}, err
		}
		f, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return []SendMessage{}, err
		}
		message := SendMessage{
			Address: record[0],
			Amount:  f,
		}
		list = append(list, message)
	}
	return
}
