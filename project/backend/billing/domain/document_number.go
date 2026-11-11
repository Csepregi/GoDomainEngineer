package domain

import (
	"fmt"
	"errors"
)

// types used by the billing domain.

type DocumentSeries struct{
	series string
}
type DocumentNumber struct{
	series DocumentSeries
	number int
}

var DocumentSeriesReceipt = DocumentSeries{"R"}

func (ds DocumentSeries) String() string {
	return ds.series
}

func (ds DocumentSeries) IsZero() bool {
	return ds.series == ""
}

func (d DocumentNumber) IsZero() bool {
	return d.series.IsZero() && d.number == 0
}

func (d DocumentNumber) String() string {
    return fmt.Sprintf("%s-%08d", d.series, d.number)
}

func NewDocumentNumber(series DocumentSeries, number int) (DocumentNumber, error) {
	if series.IsZero() {
		return DocumentNumber{}, fmt.Errorf("document series must not be empty")
	}
	if number <= 0 {
		return DocumentNumber{}, fmt.Errorf("document number must be greater than zero")
	}
	return DocumentNumber{series: series, number: number}, nil
}

func NewDocumentSeries(series string) (DocumentSeries, error) {
	if series == "" {
		return DocumentSeries{}, errors.New("document series must not be empty")
	}
	return DocumentSeries{series: series}, nil
}
