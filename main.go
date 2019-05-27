package main

import (
	ctrl "barefeed-rest/ctrl"
	driver "barefeed-rest/driver"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"log"
	"os"
)

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

// URL_SERVER=https://musing-borg-502f37.netlify.com

func main() {
	fmt.Println(os.Environ())
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set.")
	}
	client1 := os.Getenv("URL_LOCAL")
	if client1 == "" {
		log.Fatal("$URL_LOCAL must be set.")
	}
	client2 := os.Getenv("URL_SERVER")
	if client2 == "" {
		log.Fatal("$URL_SERVER must be set.")
	}
	r := gin.Default()
	db := driver.ConnectDB()

	// CORS Middleware
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{client1, client2}
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}
	r.Use(cors.New(cfg))

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
