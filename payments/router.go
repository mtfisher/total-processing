package payments

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mtfisher/total-processing/core"
)

func PaymentRoutes(router *gin.Engine) {
	router.POST("/payment", paymentPage)
	router.GET("/processed", processedPage)
}

func paymentPage(c *gin.Context) {
	loggedInInterface, _ := c.Get("is_logged_in")
	core := c.MustGet(core.CoreKey).(*core.Core)
	amount := c.PostForm("amount")
	reference := c.PostForm("reference")

	if strings.TrimSpace(amount) == "" {
		c.HTML(http.StatusBadRequest, "index.tmpl", gin.H{
			"title":        "Main website",
			"is_logged_in": loggedInInterface.(bool),
			"ErrorTitle":   "Amount is invalid",
			"ErrorMessage": "The given amount was not a valid price"})

		return;
	}

	if strings.TrimSpace(reference) == "" {
		c.HTML(http.StatusBadRequest, "index.tmpl", gin.H{
			"title":        "Main website",
			"is_logged_in": loggedInInterface.(bool),
			"ErrorTitle":   "Reference is invalid",
			"ErrorMessage": "No reference given"})

		return;
	}
	
	parsedAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.tmpl", gin.H{
			"title":        "Main website",
			"is_logged_in": loggedInInterface.(bool),
			"ErrorTitle":   "Amount is invalid",
			"ErrorMessage": "The given amount was not a valid price"})

		return;
	}

	payment, err := prepareCheckout(core.GetEntityId(), core.GetAccessKey(), core.GetCheckoutEndPoint(), parsedAmount, reference)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.tmpl", gin.H{
			"title":        "Main website",
			"is_logged_in": loggedInInterface.(bool),
			"ErrorTitle":   "Error with payment system",
			"ErrorMessage": err.Error()})

		return;
	}

	c.HTML(http.StatusOK, "payment.tmpl", gin.H{
		"title":        "Payment",
		"is_logged_in": loggedInInterface.(bool),
		"result_url": "http://localhost:8080/processed",
		"payment": payment})

}

func processedPage(c *gin.Context) {
	loggedInInterface, _ := c.Get("is_logged_in")
	core := c.MustGet(core.CoreKey).(*core.Core)
	id := c.Query("id")

	paymentResult, err := getPaymentStatus(core.GetEntityId(), core.GetAccessKey(), core.GetPaymentEndpoint(id))

	if err != nil {
		c.HTML(http.StatusBadRequest, "index.tmpl", gin.H{
			"title":        "Main website",
			"is_logged_in": loggedInInterface.(bool),
			"ErrorTitle":   "Error with payment system",
			"ErrorMessage": err.Error()})

		return;
	}

	c.HTML(http.StatusOK, "status.tmpl", gin.H{
		"title":        "Payment",
		"is_logged_in": loggedInInterface.(bool),
		"result": paymentResult})
}
