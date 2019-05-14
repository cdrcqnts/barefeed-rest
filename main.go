package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/subosito/gotenv"
	"go-mongo/aux"
	"go-mongo/con"
	"go-mongo/mdl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

var client *mongo.Client

// POST "/link"
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
		err := aux.UrlIsFeed(feed.URL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = db.InsertOne(context.TODO(), feed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": mdl.FeedClient{SID: sid, CID: cid}})
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
		err := db.FindOne(context.Background(), filter).Decode(&foundFeed)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": con.SIDNotExist})
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
		err = db.FindOne(context.Background(), filter).Decode(&foundFeed)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": con.SIDURLExist})
			return
		}
		// check if url is a valid feed address
		err = aux.UrlIsFeed(feed.URL)
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
		c.JSON(http.StatusOK, gin.H{"success": mdl.FeedClient{SID: sid, CID: cid}})
	}
}

// GET "feeds/:sid"
func GetFeeds(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO



		c.JSON(http.StatusOK, gin.H{"success": nil})
	}
}

// GET "feeds/:sid/:cid"
func GetFeed(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		cid := c.Param("cid")
		feed := mdl.Feed{SID: sid, CID: cid}
		// check if request body format is correct
		if err := c.ShouldBindJSON(&feed); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Fetch feed from url
		
		// TODO <<<

		c.JSON(http.StatusOK, gin.H{"success": nil})
	}
}

// DELETE "feeds/:sid/:cid"
func DeleteFeed(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO



		c.JSON(http.StatusOK, gin.H{"success": nil})
	}
}

// UPDATE "feeds/:sid/:cid"
func UpdateFeed(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO



		c.JSON(http.StatusOK, gin.H{"success": nil})
	}
}

func ConnectDB() *mongo.Collection {
	fmt.Println("Starting server...")
	url := os.Getenv("MONGO_URL")
	db := os.Getenv("MONGO_DB")
	col := os.Getenv("MONGO_COLLECTION")
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	res := client.Database(db).Collection(col)
	return res
}

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	r := gin.Default()
	db := ConnectDB()

	r.POST("/feeds", NewSlotNewFeed(db))      // post new URL without slot -> return slot ID
	r.POST("/feeds/:sid", OldSlotNewFeed(db)) // post new URL with slot -> slot ID

	r.GET("feeds/:sid", GetFeeds(db)) // get all channels by SID

	r.GET("feeds/:sid/:cid", GetFeed(db))       // get channel feed by SID and CID
	r.DELETE("feeds/:sid/:cid", DeleteFeed(db)) // delete channel by SID and CID
	r.PUT("feeds/:sid/:cid", UpdateFeed(db))    // update channel url by SID and CID

	err := r.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}

	//filter := bson.D{{"name", "Ash"}}

	//update := bson.D{
	//	{"$inc", bson.D{
	//		{"age", 1},
	//	}},
	//}

	//var result Trainer
	//
	//err = collection.FindOne(context.TODO(), filter).Decode(&result)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Found a single document: %+v\n", result)

	// FIND MULTIPLE ELEMENTS

	//// Pass these options to the Find method
	//findOptions := options.Find()
	//findOptions.SetLimit(2)
	//
	//// Here's an array in which you can store the decoded documents
	//var results []Trainer
	//
	//// Passing bson.D{{}} as the filter matches all documents in the collection
	//cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Finding multiple documents returns a cursor
	//// Iterating through the cursor allows us to decode documents one at a time
	//for cur.Next(context.TODO()) {
	//
	//	// create a value into which the single document can be decoded
	//	var elem Trainer
	//	err := cur.Decode(&elem)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	results = append(results, elem)
	//}
	//
	//if err := cur.Err(); err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Close the cursor once finished
	//if err := cur.Close(context.TODO()); err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	//
	//
	//err = client.Disconnect(context.TODO())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Connection to MongoDB closed.")
}
