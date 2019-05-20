package ctrl

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"go-mongo/cnt"
	"go-mongo/mdl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// GET "feeds/:sid"
func GetFeeds(db *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		feeds, err := GetURLs(db, sid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cnls, err := FeedsToChannels(feeds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": cnls})
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
		cnl, err := FeedToChannel(feed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": cnl})
	}
}

func GetURLs(db *mongo.Collection, sid string) ([]*mdl.Feed, error) {
	var res []*mdl.Feed
	filter := bson.D{{"sid", sid}}
	cur, err := db.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var f mdl.Feed
		err := cur.Decode(&f)
		if err != nil {
			return nil, err
		}
		res = append(res, &f)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	err = cur.Close(context.TODO())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FeedToChannel(feed mdl.Feed) (*mdl.Channel, error) {
	fp := gofeed.NewParser()
	feedRaw, err := fp.ParseURL(feed.URL)
	if err != nil {
		return nil, err
	}
	ps := []mdl.Podcast{}
	dayZero := time.Date(0, time.January, 1, 1, 0, 0, 0, time.UTC)
	for _, i := range feedRaw.Items {
		p := mdl.Podcast{
			Url:         "-",
			Title:       "-",
			Description: "-",
			Duration:    "-",
			Released:    dayZero,
			Image:       "",
		}
		if len(i.Title) > 0 {
			p.Title = i.Title
		}
		if len(i.Description) > 0 {
			p.Description = i.Description
		}
		if i.PublishedParsed != nil {
			p.Released = *i.PublishedParsed
		}
		if p.Released.Equal(dayZero) && i.UpdatedParsed != nil {
			p.Released = *i.UpdatedParsed
		}
		if len(i.Enclosures) > 0 {
			if len(i.Enclosures[0].URL) > 0 {
				p.Url = i.Enclosures[0].URL
			}
		}
		if i.ITunesExt != nil {
			if len(i.ITunesExt.Duration) > 0 {
				p.Duration = i.ITunesExt.Duration
			}
		}
		if i.Image != nil {
			if len(i.Image.URL) > 0 {
				p.Image = i.Image.URL
			}
		}
		ps = append(ps, p)
	}
	cnl := mdl.Channel{
		CID:         feed.CID,
		SID:         feed.SID,
		Url:         feed.URL,
		Web:         "",
		Title:       "-",
		Description: "-",
		Image:       "",
		Updated:     dayZero,
		Podcasts:    ps,
	}
	if len(feedRaw.Link) > 0 {
		cnl.Web = feedRaw.Link
	}
	if len(feedRaw.Title) > 0 {
		cnl.Title = feedRaw.Title
	}
	if len(feedRaw.Description) > 0 {
		cnl.Description = feedRaw.Description
	}
	if feedRaw.PublishedParsed != nil {
		cnl.Updated = *feedRaw.PublishedParsed
	}
	if cnl.Updated.Equal(dayZero) && feedRaw.UpdatedParsed != nil {
		cnl.Updated = *feedRaw.UpdatedParsed
	}
	if feedRaw.Image != nil {
		if len(feedRaw.Image.URL) > 0 {
			cnl.Image = feedRaw.Image.URL
		}
	}
	return &cnl, nil
}

func FeedsToChannels(feeds []*mdl.Feed) ([]*mdl.Channel, error) {
	var res []*mdl.Channel
	for _, f := range feeds {
		c, err := FeedToChannel(*f)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}
