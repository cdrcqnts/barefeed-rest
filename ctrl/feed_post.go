package ctrl

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go-mongo/aux"
	"go-mongo/cnt"
	"go-mongo/mdl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// POST "/feeds"
func NewSlotNewFeed(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := xid.New().String()
		cid := xid.New().String()
		feed := mdl.Feed{SID: sid, CID: cid}
		// check if request body format is correct
		if err := c.ShouldBindJSON(&feed); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// check if url is a valid feed address
		err := aux.UrlIsAudioFeed(feed.URL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = db.InsertOne(context.TODO(), feed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		res := struct {
			SID string `json:"sid"`
			CID string `json:"cid"`
		}{SID: feed.SID, CID: feed.CID,}
		c.JSON(http.StatusOK, gin.H{"success": res})
	}
}

// POST "/feeds/:sid"
func OldSlotNewFeed(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		cid := xid.New().String()
		feed := mdl.Feed{SID: sid, CID: cid}
		// check if request body format is correct
		if err := c.ShouldBindJSON(&feed); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// check if an entry with given SID exists
		var foundFeed mdl.Feed
		filter := bson.D{{"sid", sid}}
		err := db.FindOne(context.TODO(), filter).Decode(&foundFeed)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": cnt.ErrInvalidSID})
			return
		}
		// check if URL for given SID already exists
		filter = bson.D{
			{"$and",
				bson.A{
					bson.D{{"sid", sid}},
					bson.D{{"url", feed.URL}},
				}},
		}
		err = db.FindOne(context.TODO(), filter).Decode(&foundFeed)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": cnt.ErrUrlExists})
			return
		}
		// check if url is a valid feed address
		err = aux.UrlIsAudioFeed(feed.URL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// insert new entry with fresh CID
		_, err = db.InsertOne(context.TODO(), feed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		res := struct {
			SID string `json:"sid"`
			CID string `json:"cid"`
		}{SID: feed.SID, CID: feed.CID,}
		c.JSON(http.StatusOK, gin.H{"success": res})
	}
}
