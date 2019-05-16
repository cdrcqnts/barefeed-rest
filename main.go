package main

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"go-mongo/ctrl"
	"go-mongo/driver"
	"log"
)

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	valid.SetFieldsRequiredByDefault(true)
}

func main() {
	r := gin.Default()
	db := driver.ConnectDB()

	r.POST("/feeds", ctrl.NewSlotNewFeed(db))
	r.POST("/feeds/:sid", ctrl.OldSlotNewFeed(db))
	r.GET("feeds/:sid", ctrl.GetCIDs(db))
	r.DELETE("feeds/:sid", ctrl.DeleteSlot(db))
	r.GET("feeds/:sid/:cid", ctrl.GetFeed(db))
	r.DELETE("feeds/:sid/:cid", ctrl.DeleteFeed(db))

	err := r.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}
