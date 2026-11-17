package domain

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/shopspring/decimal"
	"project/shared/domain/shared"
)

var nonDigits = regexp.MustCompile(`\D+`)

func UnmarshalDocumentNumber(series DocumentSeries, docNumber string) (DocumentNumber, error) {
	trimmed := strings.TrimPrefix(docNumber, series.String())
	digits := nonDigits.ReplaceAllString(trimmed, "")
	number, err := strconv.Atoi(digits)
	if err != nil {
		return DocumentNumber{}, err
	}

	return DocumentNumber{
		series: series,
		number: number,
	}, nil
}

func UnmarshalTaxRate(
	rate decimal.Decimal,
	taxType TaxType,
) TaxRate {
	return TaxRate{
		rate:    rate,
		taxType: taxType,
	}
}

func UnmarshalPriceBreakdown(
	rate TaxRate,
	unitNetAmount decimal.Decimal,
	unitTaxAmount decimal.Decimal,
	unitGrossAmount decimal.Decimal,
	netAmount decimal.Decimal,
	taxAmount decimal.Decimal,
	grossAmount decimal.Decimal,
) PriceBreakdown {
	return PriceBreakdown{
		rate:            rate,
		unitNetAmount:   unitNetAmount,
		unitTaxAmount:   unitTaxAmount,
		unitGrossAmount: unitGrossAmount,
		netAmount:       netAmount,
		taxAmount:       taxAmount,
		grossAmount:     grossAmount,
	}
}

func UnmarshalLineItem(
	lineItemUUID LineItemUUID,
	name string,
	breakdown PriceBreakdown,
	quantity int,
) LineItem {
	return LineItem{
		uuid:      lineItemUUID,
		name:      name,
		breakdown: breakdown,
		quantity:  quantity,
	}
}

func UnmarshalPriceBreakdownSummary(
	netAmount decimal.Decimal,
	taxAmount decimal.Decimal,
	grossAmount decimal.Decimal,
	taxes []TaxSummary,
) PriceBreakdownSummary {
	return PriceBreakdownSummary{
		netAmount:   netAmount,
		taxAmount:   taxAmount,
		grossAmount: grossAmount,
		taxes:       taxes,
	}
}


func UnmarshalDocument(
	documentUUID DocumentUUID,
	externalReference *string,
	documentNumber DocumentNumber,
	documentType DocumentType,
	issueDate time.Time,
	currency shared.Currency,
	seller shared.LegalEntity,
	buyer shared.LegalEntity,
	items []LineItem,
	summary PriceBreakdownSummary,
) *Document {
	return &Document{
		documentUUID: DocumentUUID(documentUUID),
		externalReference: &ExternalReference(externalReference),
		documentNumber: DocumentNumber(documentNumber),
		documentType: DocumentType(documentType),
		issueDate: IssueDate(issueDate),
		currency: Currency(currency),
		seller: Company(seller),
		buyer: Company(buyer),
		lineItems: LineItems(items),
		priceBreakdownSummary: PriceBreakdownSummary(summary),
	}
}