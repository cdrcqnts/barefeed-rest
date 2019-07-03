package feeds

import (
	"context"
	"net/http"

	"github.com/cdrcqnts/barefeed-rest/api/types"
	"github.com/cdrcqnts/barefeed-rest/aux"
	"github.com/cdrcqnts/barefeed-rest/constants"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewSlotNewFeed returns the content of a feed given a feed URL
// and adds the feed URL with a fresh sid and a fresh cid to the database.
// Endpoint: POST "/feeds"
func NewSlotNewFeed(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := xid.New().String()
		cid := xid.New().String()
		link := types.Link{SID: sid, CID: cid}
		// check if request body format is correct
		if err := c.ShouldBindJSON(&link); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// check if url is a valid feed address
		err := aux.UrlIsAudioFeed(link.URL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = db.InsertOne(context.TODO(), link)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		cnl, err := aux.LinkToFeed(link)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": cnl})
	}
}

// OldSlotNewFeed returns the content of a feed given a feed URL
// and adds the feed URL with a given sid and a fresh cid to the database.
// Endpoint: POST "/feeds/:sid"
func OldSlotNewFeed(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		cid := xid.New().String()
		link := types.Link{SID: sid, CID: cid}
		// check if request body format is correct
		if err := c.ShouldBindJSON(&link); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// check if url is a valid feed address
		err := aux.UrlIsAudioFeed(link.URL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// check if an entry with given SID exists
		var foundLink types.Link
		filter := bson.D{{"sid", sid}}
		err = db.FindOne(context.TODO(), filter).Decode(&foundLink)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSID})
			return
		}
		// check if URL for given SID already exists
		filter = bson.D{
			{"$and",
				bson.A{
					bson.D{{"sid", sid}},
					bson.D{{"url", link.URL}},
				}},
		}
		err = db.FindOne(context.TODO(), filter).Decode(&foundLink)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrUrlExists})
			return
		}

		// insert new entry with fresh CID
		_, err = db.InsertOne(context.TODO(), link)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cnl, err := aux.LinkToFeed(link)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": cnl})
	}
}
