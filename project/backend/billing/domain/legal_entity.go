package domain

import (
	"errors"
	"eats/backend/common/shared"
)

type LegalEntity struct {
	name string
	taxID *shared.TaxID
	address shared.Address
	isZero bool
}

func NewLegalEntity(name string, address shared.Address, taxID *shared.TaxID) (LegalEntity, error) {
	if name == "" {
		return LegalEntity{}, errors.New("name can't be empty")
	}

	if address.IsZero() {
		return LegalEntity{}, errors.New("address can't be empty")
	}

	return LegalEntity{
		name: name,
		address: address,
		taxID: taxID,
	}, nil
}
