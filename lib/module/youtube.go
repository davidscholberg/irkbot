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
)

func ConfigYoutube(cfg *configure.Config) {
	// This is an optional function to configure the module. It is called only
	// once when irkbot starts up.
	// This function can be omitted if no configuration is needed.
}

func HelpYoutube() []string {
	s := "yt <phrase> - search youtube for the given phrase and link the top result"
	return []string{s}
}

func Youtube(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	if !strings.HasPrefix(in.Src, "#") {
		actions.Say("youtube searches not allowed in PMs")
		return
	}

	msg := strings.Join(in.MsgArgs[1:], " ")
	//fetch API key from config
	apiKey := cfg.Modules["youtube"]["api_key"]

	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating youtube client: %v", err)
		actions.Say("error creating youtube client")
		return
	}
	call := service.Search.List("id,snippet").Q(msg).MaxResults(1)
	resp, err := call.Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error performing youtube search: %v", err)
		actions.Say("error performing youtube search")
		return
	}
	var video = "no results found! ¯\\_(ツ)_/¯"
	for _, item := range resp.Items {
		switch item.Id.Kind {
		case "youtube#video":
			video = fmt.Sprintf(
				"%s - https://www.youtube.com/watch?v=%s",
				item.Snippet.Title,
				item.Id.VideoId,
			)
			break
		}
	}

	actions.Say(video)
}
