package payments

import (
	"errors"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/sonh/qs"
)

func prepareCheckout(entityId string, accessToken string, endpoint string, amount float64, merchantTransactionId string) (Payments, error) {

	checkoutForm := checkout{
		EntityId: entityId,
		Amount: amount,
		Currency: "GBP", //hard coding for speed
		PaymentType: "DB", //hard coding for speed
		MerchantTransactionId: merchantTransactionId} 
	encoder := qs.NewEncoder()
	encodedCheckoutForm, err := encoder.Values(checkoutForm)
	if err != nil {
		return Payments{}, errors.New("could not encode form")
	}

	client := &http.Client{}
    r, err := http.NewRequest("POST", endpoint, strings.NewReader(encodedCheckoutForm.Encode()))
    if err != nil {
        return Payments{}, err
    }
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(encodedCheckoutForm.Encode())))
	r.Header.Add("Authorization", "Bearer " + accessToken)

    res, err := client.Do(r)
	if err != nil {
        return Payments{}, err
    }
	defer res.Body.Close()
	bodyCon, err := ioutil.ReadAll(res.Body)
	if err != nil {
        return Payments{}, err
    }

	jsonRes := checkoutResponse{}
    err = json.Unmarshal(bodyCon, &jsonRes)
	if err != nil {
        return Payments{}, err
    }

	return Payments{
		CheckoutID: jsonRes.Id,
		MerchantTransactionId: merchantTransactionId,
		Result: jsonRes.Result.Code,
		Description: jsonRes.Result.Description}, nil
}

func getPaymentStatus(entityId string, accessToken string, endpoint string) (PaymentResponse, error) {
	client := &http.Client{}
    r, err := http.NewRequest("GET", endpoint + "?entityId=" + entityId, nil)
    if err != nil {
        return PaymentResponse{}, err
    }
	r.Header.Add("Authorization", "Bearer " + accessToken)

    res, err := client.Do(r)
	if err != nil {
        return PaymentResponse{}, err
    }
	defer res.Body.Close()
	bodyCon, err := ioutil.ReadAll(res.Body)

	jsonRes := PaymentResponse{}
    err = json.Unmarshal(bodyCon, &jsonRes)
	if err != nil {
        return PaymentResponse{}, err
    }

	return jsonRes, nil
}