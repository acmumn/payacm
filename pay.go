package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

type Payment struct {
	Amount uint64 `binding:"required" json:"amount"`
	Email  string `binding:"required" json:"email"`
	Reason string `binding:"required" json:"reason"`
	Token  string `binding:"required" json:"token"`
}

func pay(c *gin.Context) {
	var payment Payment
	if err := c.BindJSON(&payment); err != nil {
		log.Println("Error binding in payment", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if payment.Amount <= 0 || payment.Email == "" || payment.Reason == "" || payment.Token == "" {
		c.JSON(http.StatusBadRequest, payment)
		return
	}

	chargeParams := &stripe.ChargeParams{
		Amount:   payment.Amount,
		Currency: "usd",
		Desc:     payment.Reason,
	}
	chargeParams.SetSource(payment.Token)

	_, err := charge.New(chargeParams)
	if err != nil {
		log.Println("Error charging card", err)
		c.JSON(http.StatusServiceUnavailable, err)
		return
	}

	err = mail(payment)
	if err != nil {
		log.Println("Error sending mail", err)
		c.JSON(http.StatusBadGateway, err)
		return
	}

	mail_template.Execute(c.Writer, payment)
}
