// main.go
package main

import (
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mtfisher/total-processing/core"
	"github.com/mtfisher/total-processing/users"
	"github.com/mtfisher/total-processing/payments"
	"github.com/patrickmn/go-cache"
)

func main() {

	db, err := bolt.Open("tp.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(users.UserBucket))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	c := cache.New(5*time.Hour, 10*time.Minute)

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	coreApp := core.Initialize(db, c, os.Getenv("ENTITY_ID"), os.Getenv("ACCESS_TOKEN"), os.Getenv("API_BASE_URL"))

	// Set the router as the default one provided by Gin
	router := gin.Default()
	router.Use(coreApp.SetCore())
	router.Use(users.SetUserStatus())

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/*")

	// Initialize the router
	users.AuthedRoutes(router)
	users.LoginRoutes(router)
	payments.PaymentRoutes(router)

	// Start serving the application
	router.Run(":8080")
}
