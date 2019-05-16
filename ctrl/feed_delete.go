package ctrl

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func DeleteSlot(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		// check if URL for given SID already exists
		filter := bson.D{{"sid", sid}}
		res, err := db.DeleteMany(context.TODO(), filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": res})
	}
}

// DELETE "feeds/:sid/:cid"
func DeleteFeed(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		cid := c.Param("cid")
		filter := bson.D{
			{"$and",
				bson.A{
					bson.D{{"sid", sid}},
					bson.D{{"cid", cid}},
				}},
		}
		_, err := db.DeleteOne(context.TODO(), filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprint("Deleted entry with SID %x, CID %y", sid, cid)})
	}
}