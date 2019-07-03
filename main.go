package main

// FIXME: prefix with github
import (
	"log"
	"os"

	"github.com/cdrcqnts/barefeed-rest/ctrl"
	"github.com/cdrcqnts/barefeed-rest/driver"

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
	db := driver.ConnectDB()

	// CORS Middleware
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{client}
	r.Use(cors.New(corsCfg))

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
