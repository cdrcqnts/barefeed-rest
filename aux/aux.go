package aux

import (
	"errors"
	"github.com/mmcdole/gofeed"
	"go-mongo/cnt"
	"strings"
)

func UrlIsAudioFeed(url string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return errors.New(cnt.ErrParseURL)
	}
	if len(feed.Items) > 0 && len(feed.Items[0].Enclosures) > 0 && len(feed.Items[0].Enclosures[0].URL) > 0  {
		for _, af := range cnt.AUDIO_FORMATS {
			if strings.HasSuffix(feed.Items[0].Enclosures[0].URL, af) {
				return nil
			}
		}
	}
	return errors.New(cnt.ErrNoAudioFile)
}