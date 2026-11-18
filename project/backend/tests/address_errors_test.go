// This file contains tests that are executed to verify your solution.
// It's read-only, so all modifications will be ignored.
package tests_test

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ordersclient "eats/backend/orders/api/http/client"
)

func TestComponent_RegisterCustomer_InvalidAddressReturnsAllDetails(t *testing.T) {
	t.Parallel()
	clients := newTestClients(t)

	ctx := t.Context()

	resp, err := clients.Orders.RegisterCustomerWithResponse(ctx, ordersclient.RegisterCustomer{
		Name:        "Invalid Address Customer",
		Email:       "invalid-address@example.com",
		Address:     ordersclient.Address{},
		PhoneNumber: "+1 555 0100",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	require.NotNil(t, resp.JSON400)

	expected := ordersclient.ErrorResponse{
		Slug: "invalid-address",
		Details: []ordersclient.ErrorDetails{
			{ErrorSlug: "address-line1-required"},
			{ErrorSlug: "address-postal-code-required"},
			{ErrorSlug: "address-city-required"},
			{ErrorSlug: "address-country-code-required"},
		},
	}

	// Messages, entity types, and detail order come from the student's NewAddress.
	// Later exercises run this test against their code, so only slugs are checked.
	diff := cmp.Diff(expected, *resp.JSON400,
		cmpopts.IgnoreFields(ordersclient.ErrorResponse{}, "Message"),
		cmpopts.IgnoreFields(ordersclient.ErrorDetails{}, "Message", "EntityType", "EntityId"),
		cmpopts.SortSlices(func(a, b ordersclient.ErrorDetails) bool { return a.ErrorSlug < b.ErrorSlug }),
	)
	assert.Empty(t, diff)
}
