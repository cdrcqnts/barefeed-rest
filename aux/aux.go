package aux

import (
	"errors"
	"github.com/mmcdole/gofeed"
	"poasy-rest/cnt"
	"strings"
)

func UrlIsAudioFeed(url string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return errors.New(cnt.ErrParseURL)
	}
	i := feed.Items
	e := feed.Items[0].Enclosures
	u := feed.Items[0].Enclosures[0].URL
	if len(i) > 0 && len(e) > 0 && len(u) > 0  {
		for _, af := range cnt.AUDIO_FORMATS {
			if strings.HasSuffix(u, af) {
				return nil
			}
		}
	}
	return errors.New(cnt.ErrNoAudioFile)
}