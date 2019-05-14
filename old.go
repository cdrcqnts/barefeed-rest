package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// DONE
func NewSlot(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		slot := Link{SID: xid.New().String(), CID: xid.New().String()}
		res, err := db.InsertOne(context.TODO(), slot)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": res.InsertedID})
	}
}

//TODO return parsed Channels with Podcasts ?
func GetSlot(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var slot Slot
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.D{{"_id", id}}
		err = db.FindOne(context.Background(), filter).Decode(&slot)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": slot})
	}
}

// DONE
func DeleteSlot(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var slot Slot
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.D{{"_id", id}}
		err = db.FindOneAndDelete(context.Background(), filter).Decode(&slot)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprint("Delete slot with id: %x", id)})
	}
}

//func GetFeed(db *mongo.Collection) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// TODO return feed for channel and its podcasts
//		var link Link
//		var slot SlotR
//		if err := c.ShouldBindJSON(&link); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//		id, err := primitive.ObjectIDFromHex(c.Param("id"))
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//		filter := bson.D{{"_id", id}}
//		err = db.FindOne(context.Background(), filter).Decode(&slot)
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//
//		fp := gofeed.NewParser()
//		feed, err := fp.ParseURL(link.URL)
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse URL. Please make sure the feed URL is valid."})
//			return
//		}
//		cnls = append(cnls, FeedToStruct(cnl.ID, cnl.Slot, cnl.Url, feed))
//
//		c.JSON(http.StatusOK, gin.H{"success": slot})
//	}
//}

func AddLink(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var slot Slot
		var link Link
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err = c.ShouldBindJSON(&link); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.D{{"_id", id}}
		update := bson.D{
			{"$addToSet", bson.D{
				{"links", link},
			}},
		}
		err = db.FindOneAndUpdate(context.Background(), filter, update).Decode(&slot)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// second request to return updated list of links
		err = db.FindOne(context.Background(), filter).Decode(&slot)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": slot})
	}
}

func DeleteLink(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var slot Slot
		var link Link
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err = c.ShouldBindJSON(&link); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.D{{"_id", id}}
		update := bson.D{
			{"$pull", bson.D{
				{"links", link},
			}},
		}
		err = db.FindOneAndUpdate(context.Background(), filter, update).Decode(&slot)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// second request to return updated list of links
		err = db.FindOne(context.Background(), filter).Decode(&slot)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": slot})
	}
}
