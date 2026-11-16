package domain

import "github.com/shopspring/decimal"

type PriceBreakdownSummary struct {
	netAmount   decimal.Decimal
	taxAmount   decimal.Decimal
	grossAmount decimal.Decimal
	taxes       []TaxSummary
}

type taxRateKey struct {
	rate    string
	taxType string
}

func (t TaxRate) key() taxRateKey {
	return taxRateKey{
		rate:    t.rate.String(),
		taxType: t.taxType.String(),
	}
}



func summarizeLineItems(lineItems []LineItem) PriceBreakdownSummary {
	netAmount := decimal.Zero
	taxAmount := decimal.Zero
	grossAmount := decimal.Zero
	

	for _, lineItem := range lineItems {
		breakdown := lineItem.PriceBreakdown()
		netAmount = netAmount.Add(breakdown.NetAmount())
		taxAmount = taxAmount.Add(breakdown.TaxAmount())
		grossAmount = grossAmount.Add(breakdown.GrossAmount())
	}

	return PriceBreakdownSummary{
		netAmount:   netAmount,
		taxAmount:   taxAmount,
		grossAmount: grossAmount,
		taxes:       newGroupedTaxes(lineItems),
	}
}

func (p PriceBreakdownSummary) NetAmount() decimal.Decimal {
	return p.netAmount
}

func (p PriceBreakdownSummary) TaxAmount() decimal.Decimal {
	return p.taxAmount
}

func (p PriceBreakdownSummary) GrossAmount() decimal.Decimal {
	return p.grossAmount
}

func (p PriceBreakdownSummary) Taxes() []TaxSummary {
	return p.taxes
}

type TaxSummary struct {
	taxRate   TaxRate
	netAmount decimal.Decimal
	taxAmount decimal.Decimal
}

func (t TaxSummary) TaxRate() TaxRate {
	return t.taxRate
}

func (t TaxSummary) NetAmount() decimal.Decimal {
	return t.netAmount
}

func (t TaxSummary) TaxAmount() decimal.Decimal {
	return t.taxAmount
}


func newTaxSummary(priceBreakdown PriceBreakdown) *TaxSummary {
	return &TaxSummary{
		taxRate:   priceBreakdown.TaxRate(),
		netAmount: priceBreakdown.NetAmount(),
		taxAmount: priceBreakdown.TaxAmount(),
	}
}

func (t *TaxSummary) add(priceBreakdown PriceBreakdown) {
	t.netAmount = t.netAmount.Add(priceBreakdown.NetAmount())
	t.taxAmount = t.taxAmount.Add(priceBreakdown.TaxAmount())
}

func newGroupedTaxes(lineItems []LineItem) []TaxSummary {
	taxes := map[taxRateKey]*TaxSummary{}

	for _, l := range lineItems {
		price := l.PriceBreakdown()
		key := price.rate.key()

		summary, ok := taxes[key]
		if ok {
			summary.add(price)
		} else {
			taxes[key] = newTaxSummary(price)
		}
	}

	keys := make([]taxRateKey, 0, len(taxes))
	for k := range taxes {
		keys = append(keys, k)
	}

	result := make([]TaxSummary, 0, len(taxes))
	for _, k := range keys {
		result = append(result, *taxes[k])
	}

	return result
}