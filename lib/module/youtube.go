package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"os"
	"strings"
	"time"
)

func helpYoutubeSearch() []string {
	s := "yt <phrase> - search youtube for the given phrase and link the top result"
	return []string{s}
}

func youtubeSearch(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	if !strings.HasPrefix(in.Src, "#") {
		actions.say("youtube searches not allowed in PMs")
		return
	}

	msg := strings.Join(in.MsgArgs[1:], " ")
	//fetch API key from config
	apiKey := cfg.Modules["youtube"]["api_key"]

	client := &http.Client{
		Timeout:   time.Duration(cfg.Http.Timeout) * time.Second,
		Transport: &transport.APIKey{Key: apiKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating youtube client: %v\n", err)
		actions.say("error creating youtube client")
		return
	}
	call := service.Search.List("id,snippet").Q(msg).MaxResults(1)
	resp, err := call.Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error performing youtube search: %v\n", err)
		actions.say("error performing youtube search")
		return
	}
	var result = "no results found! ¯\\_(ツ)_/¯"
	for _, item := range resp.Items {
		switch item.Id.Kind {
		case "youtube#video":
			result = fmt.Sprintf(
				"%s - https://www.youtube.com/watch?v=%s",
				item.Snippet.Title,
				item.Id.VideoId,
			)
		case "youtube#channel":
			result = fmt.Sprintf(
				"%s - https://www.youtube.com/channel/%s",
				item.Snippet.Title,
				item.Id.ChannelId,
			)
		case "youtube#playlist":
			playlistId, err := getExternalPlaylistId(service, item.Id.PlaylistId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting playlist id: %v\n", err)
				break
			}
			result = fmt.Sprintf(
				"%s - https://www.youtube.com/playlist?list=%s",
				item.Snippet.Title,
				playlistId,
			)
		}
	}

	actions.say(result)
}

// getExternalPlaylistId does an extra call to the youtube api to get the
// actual external playlist ID that can be used in URLs.
func getExternalPlaylistId(service *youtube.Service, playlistId string) (string, error) {
	call := service.Playlists.List("id").Id(playlistId).MaxResults(1)
	resp, err := call.Do()
	if err != nil {
		return "", err
	}
	for _, item := range resp.Items {
		return item.Id, nil
	}
	return "", nil
}
