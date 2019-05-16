package aux

import (
	"errors"
	"github.com/mmcdole/gofeed"
	"go-mongo/cnt"
	"go-mongo/mdl"
	"strings"
	"time"
)

func UrlIsAudioFeed(url string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return errors.New(cnt.ErrParseURL)
	}
	// verify that feed contains a link to an audio file
	if len(feed.Items) > 0 && len(feed.Items[0].Enclosures) > 0 && len(feed.Items[0].Enclosures[0].URL) > 0  {
		for _, af := range cnt.AUDIO_FORMATS {
			if strings.HasSuffix(feed.Items[0].Enclosures[0].URL, af) {
				return nil
			}
		}
	}
	return errors.New(cnt.ErrNoAudioFile)
}

func SlimFeed(feed mdl.Feed, feedRaw *gofeed.Feed) mdl.Channel {
	ps := []mdl.Podcast{}
	dayZero := time.Date(0, time.January, 1, 1, 0, 0, 0, time.UTC)
	for _, i := range feedRaw.Items {
		p := mdl.Podcast{
			Url:         "",
			Title:       "",
			Description: "",
			Duration:    "",
			Released:    dayZero,
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
				p.Duration = i.Enclosures[0].Length
			}
		}
		ps = append(ps, p)
	}
	cnl := mdl.Channel{
		CID:         feed.CID,
		SID:         feed.SID,
		Url:         feed.URL,
		Web:         "",
		Title:       "",
		Description: "",
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
	return cnl
}
