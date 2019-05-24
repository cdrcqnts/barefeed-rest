package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"barefeed-rest/ctrl"
	"barefeed-rest/driver"
	"log"
)

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	r := gin.Default()
	db := driver.ConnectDB()

	// CORS Middleware
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{"http://localhost:1234"}
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}
	r.Use(cors.New(cfg))

	r.POST("/feeds", ctrl.NewSlotNewFeed(db))
	r.POST("/feeds/:sid", ctrl.OldSlotNewFeed(db))
	r.GET("feeds/:sid", ctrl.GetFeeds(db))
	r.DELETE("feeds/:sid", ctrl.DeleteSlot(db))
	r.GET("feeds/:sid/:cid", ctrl.GetFeed(db))
	r.DELETE("feeds/:sid/:cid", ctrl.DeleteFeed(db))

	err := r.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}
