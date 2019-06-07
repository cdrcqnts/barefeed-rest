package aux

import (
	cnt "barefeed-rest/cnt"
	mdl "barefeed-rest/mdl"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UrlIsAudioFeed returns no error if the passed string refers to a feed containing an enclosed audio file
func UrlIsAudioFeed(url string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return errors.New(cnt.ErrParseURL)
	}
	i := feed.Items
	e := feed.Items[0].Enclosures
	u := feed.Items[0].Enclosures[0].URL
	if len(i) > 0 && len(e) > 0 && len(u) > 0 {
		for _, af := range cnt.AUDIO_FORMATS {
			if strings.HasSuffix(u, af) {
				return nil
			}
		}
	}
	return errors.New(cnt.ErrNoAudioFile)
}

// GetURLs returns a list of all feeds for a given sid
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
	if len(res) > 0 {
		return res, nil
	}
	return nil, errors.New(fmt.Sprint("Could not find slot with id: ", sid))
}

// FeedToChannel returns the content of a feed given a feed URL
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
			PID:         xid.New().String(),
			Url:         "-",
			Title:       "-",
			Description: "-",
			Duration:    "-",
			Released:    dayZero,
			Image:       "",
			Size:        0,
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
			if len(i.Enclosures[0].Length) > 0 {
				i, _ := strconv.Atoi(i.Enclosures[0].Length)
				p.Size = i
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

// FeedsToChannels returns the content of a list of feeds given a list of feed URLs
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
