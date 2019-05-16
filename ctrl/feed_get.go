package ctrl

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"go-mongo/aux"
	"go-mongo/cnt"
	"go-mongo/mdl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// GET "feeds/:sid"
func GetCIDs(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		var res []*string
		filter := bson.D{{"sid", sid}}
		cur, err := db.Find(context.TODO(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for cur.Next(context.TODO()) {
			var f mdl.Feed
			err := cur.Decode(&f)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			res = append(res, &f.CID)
		}
		if err := cur.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = cur.Close(context.TODO())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if res == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": cnt.ErrInvalidSID})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": res})
	}
}

// GET "feeds/:sid/:cid"
func GetFeed(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		cid := c.Param("cid")
		feed := mdl.Feed{SID: sid, CID: cid}
		filter := bson.D{
			{"$and",
				bson.A{
					bson.D{{"sid", sid}},
					bson.D{{"cid", cid}},
				}},
		}
		err := db.FindOne(context.TODO(), filter).Decode(&feed)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": cnt.ErrSIDCIDNotFound})
			return
		}
		fp := gofeed.NewParser()
		feedRaw, err := fp.ParseURL(feed.URL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": cnt.ErrParseURL})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": aux.SlimFeed(feed, feedRaw)})
	}
}
