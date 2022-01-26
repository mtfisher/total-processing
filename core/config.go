package core

import (
	"strings"
)

const checkoutUri = "v1/checkouts"
const paymentUri = "v1/checkouts/_id_/payment"

func (c Core) GetEntityId() string {
	return c.entityId;
}

func (c Core) GetAccessKey() string {
	return c.accessToken
}

func (c Core) GetCheckoutEndPoint() string {
	return c.apiBaseUrl + checkoutUri
}

func (c Core) GetPaymentEndpoint(id string) string {
	return strings.Replace(c.apiBaseUrl + paymentUri, "_id_", id, -1)
}