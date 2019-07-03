package main

import (
	"log"
	"os"

	"github.com/cdrcqnts/barefeed-rest/api/v1/feeds"
	"github.com/cdrcqnts/barefeed-rest/api/v1/ping"
	"github.com/cdrcqnts/barefeed-rest/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

// init loads the environment variables from .env before main() is  executed

// TODO remove, not needed!
func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

// main is the application entry point
// run `go run main.go` to start the server
func main() {
	// FIXME: USE YAML, OCCAMY for config
	// fmt.Println(os.Environ())
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set.")
	}
	client := os.Getenv("CLIENT")
	if client == "" {
		log.Fatal("$CLIENT must be set.")
	}
	r := gin.Default()
	// FIXME: turn into global variable
	db := database.Connect()

	// CORS Middleware
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{client}
	r.Use(cors.New(corsCfg))

	// TODO: Middeware for request limit

	v1 := r.Group("/v1")
	{
		v1.GET("/ping", ping.Ping())
		v1.POST("/feeds", feeds.NewSlotNewFeed(db))
		v1.POST("/feeds/:sid", feeds.OldSlotNewFeed(db))
		v1.GET("feeds/:sid", feeds.GetFeeds(db))
		v1.DELETE("feeds/:sid", feeds.DeleteSlot(db))
		v1.GET("feeds/:sid/:cid", feeds.GetFeed(db))
		v1.DELETE("feeds/:sid/:cid", feeds.DeleteFeed(db))
	}

	err := r.Run(":" + port)
	if err != nil {
		log.Fatalln(err)
	}
}
