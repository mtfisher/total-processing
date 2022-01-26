package payments

type Payments struct {
	CheckoutID string
	MerchantTransactionId string
	Result string
	Description string
}

type checkout struct {
	EntityId string `qs:"entityId"`
	Amount float64 `qs:"amount"`
	Currency string `qs:"currency"`
	PaymentType string `qs:"paymentType"`
	MerchantTransactionId string `qs:"merchantTransactionId"`
}

type checkoutResponse struct {
	Id string `json:"id"`
	Result TpResult `json:"result"`
}

type PaymentResponse struct {
	Id string `json:"id"`
	Result TpResult `json:"result"`
	Amount string `json:"amount"`
	TimeStamp string `json:"timestamp"`
	MerchantTransactionId string `json:"merchantTransactionId"`
}

type TpResult struct {
	Code string `json:"code"`
	Description string `json:"description"`
}

