package domain

import (
	"errors"
	"time"

	"eats/backend/common"
	"eats/backend/common/shared"
)

type DocumentType struct {
	common.Enum[DocumentTypeValues]
}

type DocumentTypeValues string

func (DocumentTypeValues) Values() []string {
	return []string{"receipt"}
}

type Document struct {
	uuid DocumentUUID
	externalReference *string
	documentType DocumentType
	documentNumber DocumentNumber
	issueDate time.Time
	currency shared.Currency
	seller LegalEntity
	buyer LegalEntity
	lineItems []LineItem
}

type LineItem struct {
	uuidID LineItemUUID
	name string
	quantity int
}




type NewDocumentData struct {
	ExternalReference *string
	IssueDate         time.Time
	Currency          shared.Currency
	Seller            LegalEntity
	Buyer             LegalEntity
	LineItems         []NewLineItemData
}

type NewLineItemData struct {
	Name     string
	Quantity int
}

var DocumentTypeReceipt = common.MustEnum[DocumentType]("receipt")

func NewReceipt(data NewDocumentData, docNumber DocumentNumber) (*Document, error) {
	if data.Buyer.IsZero() {
	return nil, errors.New("buyer can't be empty")
}
	if data.Seller.IsZero() {
		return nil, errors.New("seller can't be empty")
	}
	if shared.TaxID && data.Seller.TaxID() != nil {
		return nil, errors.New("seller must have a tax ID")
	}

	if shared.TaxID && !data.Buyer.TaxID().IsZero() {
		return nil, errors.New("receipts cannot be issued to buyers with a tax ID")
	}

	if data.Currency.IsZero() {
		return nil, errors.New("currency can't be empty")
	}	

	if data.IssueDate.IsZero() {
		return nil, errors.New("issue date can't be empty")
	}

	if data.IssueDate.After(time.Now()) {
		return nil, errors.New("issue date cannot be in the future")
	}

	if len(data.lineItems) == 0 {
		return nil, errors.New("at least one line item")
	}

	for _, lineItem := range data.lineItems {
		if lineItem.name == "" {
			return nil, errors.New("name can't be empty")
		}
		if lineItem.quantity <= 0 {
			return nil, errors.New("quantity must be positive")
		}
	}

	return &Document{
		uuid: data.uuid,
	}
}

type DocumentUUID struct {
	common.UUID
}

type LineItemUUID struct {
	common.UUID
}


func uuid() DocumentUUID {
	return DocumentUUID{
		common.UUID: common.NewUUIDv7(),
	}
}

func externalReference() *string {
	return nil
}

func documentType() DocumentType {
	return DocumentType{
		common.Enum[DocumentTypeValues]: DocumentTypeReceipt,
	}
}

func documentNumber() DocumentNumber {
	return DocumentNumber{
		number: "1234567890",
	}
}

func issueDate() time.Time {
	return time.Now()
}

func currency() shared.Currency {
	return shared.Currency{
		code: "USD",
	}
}

func seller() LegalEntity {
	return LegalEntity{
		name: "Test Seller",
		address: shared.Address{
			street: "123 Main St",
			city: "Anytown",
			state: "CA",
			zip: "12345",
		},
		taxID: &shared.TaxID{
			taxID: "1234567890",
		},
	}
}

func name() string {
	return "Test Name"
}

func quantity() int {
	return 1
}