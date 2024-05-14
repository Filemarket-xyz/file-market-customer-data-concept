package domain

import (
	"strings"

	"github.com/gocarina/gocsv"
)

type Dataset struct {
	Id       int64 `csv:"id"`
	ClientId int64 `csv:"client_id"`
	// TODO
}

func DatasetsToStrings(d []*Dataset) ([][]string, error) {
	out, err := gocsv.MarshalString(d)
	if err != nil {
		return nil, err
	}

	rows := strings.Fields(out)
	res := make([][]string, len(rows))

	for i, row := range rows {
		res[i] = strings.Split(row, ",")
	}

	return res, nil
}
