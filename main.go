package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Using default port of 3000")
		port = "3000"
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		log.Fatal("Did not provide STRIPE_SECRET_KEY")
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.LoadHTMLGlob("templates/*")
	r.Static("static", "static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("/", pay)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
