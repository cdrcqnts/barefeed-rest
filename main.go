package main

import (
	ctrl "barefeed-rest/ctrl"
	driver "barefeed-rest/driver"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

// init loads the environment variables from .env before main() is  executed
func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

// main is the application entry point
// run `go run main.go` to start the server
func main() {
	// fmt.Println(os.Environ())
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set.")
	}
	client := os.Getenv("CLIENT")
	if client == "" {
		log.Fatal("$CLIENT must be set.")
	}
	clientLocal := os.Getenv("CLIENT_LOCAL")
	if clientLocal == "" {
		log.Fatal("$CLIENT_LOCAL must be set.")
	}
	r := gin.Default()
	db := driver.ConnectDB()

	// CORS Middleware
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{client, clientLocal}
	r.Use(cors.New(cfg))

	// TODO: Middeware for request limit

	r.POST("/feeds", ctrl.NewSlotNewFeed(db))
	r.POST("/feeds/:sid", ctrl.OldSlotNewFeed(db))
	r.GET("feeds/:sid", ctrl.GetFeeds(db))
	r.DELETE("feeds/:sid", ctrl.DeleteSlot(db))
	r.GET("feeds/:sid/:cid", ctrl.GetFeed(db))
	r.DELETE("feeds/:sid/:cid", ctrl.DeleteFeed(db))

	err := r.Run(":" + port)
	if err != nil {
		log.Fatalln(err)
	}
}
