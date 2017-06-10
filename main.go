package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go"
)

func main() {
	// Get port to serve on.
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Using default port of 3000")
		port = "3000"
	}

	// Get the stripe key.
	stripe.Key = getenv("STRIPE_SECRET_KEY")

	// Fail out early if mailer variables are not defined.
	getenv("SMTP_FROM")
	getenv("SMTP_HOST")
	getenv("SMTP_PASS")
	getenv("SMTP_PORT")
	getenv("SMTP_USER")

	// Create the router.
	r := gin.New()
	r.Use(gin.Logger())

	// Connect callbacks.
	r.StaticFile("/", "static/index.html")
	r.StaticFile("/static/main.js", "static/main.js")
	r.POST("/", pay)

	// Serve.
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func getenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal(fmt.Sprintf("Did not provide %s", key))
	}
	return val
}
