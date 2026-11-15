package domain

import (
	"errors"
	"time"
	"fmt"
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
	if data.Buyer.TaxID() != nil {
		return nil, errors.New("receipts cannot be issued to buyers with a tax ID")
	}

	if data.Seller.TaxID() == nil {
		return nil, errors.New("seller must have a tax ID to issue billing documents")
	}

	if data.Currency.IsZero() {
		return nil, errors.New("currency can't be empty")
	}	

	if data.IssueDate.IsZero() {
		return nil, errors.New("issue date can't be empty")
	}

	if data.IssueDate.After(time.Now()) {
		return nil, errors.New("issue date can't be in the future.")
	}

	if len(data.LineItems) == 0 {
		return nil, errors.New("document must have at least one line item")
	}


	lineItems := make([]LineItem, 0, len(data.LineItems))
	for _, lid := range data.LineItems {
		lineItem, err := newLineItem(lid)
		if err != nil {
			return nil, fmt.Errorf("failed to create line item for document: %w", err)
		}
		lineItems = append(lineItems, lineItem)
	}

	return &Document{
		uuid: DocumentUUID{common.NewUUIDv7()},
		externalReference: data.ExternalReference,
		documentType: DocumentTypeReceipt,
		documentNumber: docNumber,
		issueDate: data.IssueDate,
		currency: data.Currency,
		seller: data.Seller,
		buyer: data.Buyer,
		lineItems: lineItems,
	}, nil
}

type DocumentUUID struct {
	common.UUID
}

type LineItemUUID struct {
	common.UUID
}

type LineItem struct {
	uuidID LineItemUUID
	name string
	quantity int
}



func newLineItem(data NewLineItemData) (LineItem, error) {
	if data.Name == "" {
		return LineItem{}, errors.New("name can't be empty")
	}

	if data.Quantity < 1 {
		return LineItem{}, errors.New("quantity must be positive")
	}

	return LineItem{
		uuidID:   LineItemUUID{common.NewUUIDv7()},
		name:     data.Name,
		quantity: data.Quantity,
	}, nil
}

func (d *Document) UUID() DocumentUUID {
	return d.uuid
}

func (d *Document) ExternalReference() *string {
	return d.externalReference
}

func (d *Document) DocumentType() DocumentType {
	return d.documentType
}

func (d *Document) DocumentNumber() DocumentNumber {
	return d.documentNumber
}

func (d *Document) Seller() LegalEntity {
	return d.seller
}

func (d *Document) Buyer() LegalEntity {
	return d.buyer
}

func (d *Document) LineItems() []LineItem {
	return d.lineItems
}

func (d *Document) IssueDate() time.Time {
	return d.issueDate
}

func (d *Document) Currency() shared.Currency {
	return d.currency
}

func (l *LineItem) Name() string {
	return l.name
}

func (l *LineItem) Quantity() int {
	return l.quantity
}

func (l *LineItem) UUID() LineItemUUID {
	return l.uuidID
}