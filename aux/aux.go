package aux

import (
	"errors"
	"github.com/mmcdole/gofeed"
	"go-mongo/mdl"
	"strings"
)

func UrlIsFeed(url string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return errors.New("Failed parsing URL. Make sure the URL provides a valid feed.")
	}
	// verify that feed contains a link to a mp3 audio file
	if !strings.HasSuffix(feed.Items[0].Enclosures[0].URL, "mp3") {
		return errors.New("It seems this feed does not provide mp3 audio files.")
	}
	return nil
}

func FeedToStruct(cid string, sid string, url string, feed *gofeed.Feed) mdl.Channel {
	ps := []mdl.Podcast{}
	for _, i := range feed.Items {
		p := mdl.Podcast{
			Url:         i.Enclosures[0].URL,
			Title:       i.Title,
			Description: i.Description,
			Duration:    i.Enclosures[0].Length,
			Released:    *i.PublishedParsed,
		}
		ps = append(ps, p)
	}
	return mdl.Channel{
		CID:         cid,
		SID:         sid,
		Url:         url,
		Web:         feed.Link,
		Title:       feed.Title,
		Description: feed.Description,
		Updated:     *feed.UpdatedParsed,
		Podcasts:    ps,
	}
}


