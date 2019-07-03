package main

// FIXME: prefix with github
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

	r.GET("/ping", ping.Ping())
	r.POST("/feeds", feeds.NewSlotNewFeed(db))
	r.POST("/feeds/:sid", feeds.OldSlotNewFeed(db))
	r.GET("feeds/:sid", feeds.GetFeeds(db))
	r.DELETE("feeds/:sid", feeds.DeleteSlot(db))
	r.GET("feeds/:sid/:cid", feeds.GetFeed(db))
	r.DELETE("feeds/:sid/:cid", feeds.DeleteFeed(db))

	err := r.Run(":" + port)
	if err != nil {
		log.Fatalln(err)
	}
}
