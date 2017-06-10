package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

type Payment struct {
	amount uint64 `binding:"required"`
	email  string `binding:"required"`
	reason string `binding:"required"`
	token  string `binding:"required"`
}

func pay(c *gin.Context) {
	var payment Payment
	c.BindJSON(&payment)

	chargeParams := &stripe.ChargeParams{
		Amount:   payment.amount,
		Currency: "usd",
		Desc:     payment.reason,
	}
	chargeParams.SetSource(payment.token)

	_, err := charge.New(chargeParams)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusServiceUnavailable, err)
		return
	}

	err = mail(payment)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, payment)
}
