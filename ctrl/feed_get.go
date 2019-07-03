package ctrl

import (
	"context"
	"net/http"

	"github.com/cdrcqnts/barefeed-rest/aux"
	"github.com/cdrcqnts/barefeed-rest/cnt"
	"github.com/cdrcqnts/barefeed-rest/mdl"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetFeeds returns the content of all feeds corresponding to the sid specified in the URL parameter
// Endpoint: GET "feeds/:sid"
func GetFeeds(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		links, err := aux.GetLinks(db, sid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cnls, err := aux.LinksToFeeds(links)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": cnls})
	}
}

// GetFeed returns the content of the feed corresponding to the sid and the cid specified in the URL parameter
// Endpoint: GET "feeds/:sid/:cid"
func GetFeed(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		cid := c.Param("cid")
		feed := mdl.Link{SID: sid, CID: cid}
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
		cnl, err := aux.LinkToFeed(feed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": cnl})
	}
}
