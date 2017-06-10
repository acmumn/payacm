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
	log.Println("a", payment)
	if err := c.BindJSON(&payment); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	log.Println("b", payment)

	if payment.amount <= 0 || payment.email == "" || payment.reason == "" || payment.token == "" {
		c.JSON(http.StatusBadRequest, payment)
		return
	}

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
