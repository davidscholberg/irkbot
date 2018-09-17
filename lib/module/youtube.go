package module

import (
        "github.com/davidscholberg/irkbot/lib/configure"
        "github.com/davidscholberg/irkbot/lib/message"
        "strings"
        "fmt"
        "log"
        "net/http"
        "google.golang.org/api/googleapi/transport"
        "google.golang.org/api/youtube/v3"
)

func ConfigYoutube(cfg *configure.Config) {
        // This is an optional function to configure the module. It is called only
is calle// once when irkbot starts up.
        // This function can be omitted if no configuration is needed.
}

func HelpYoutube() []string {
        s := "yt <phrase> - search youtube for the given phrase and link the top result"
        return []string{s}
}

func Youtube(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
        msg := strings.Join(in.MsgArgs[1:], " ")
        //fetch API key from config
        apiKey := cfg.Modules["youtube"]["api_key"]

        client := &http.Client{
                Transport: &transport.APIKey{Key: apiKey},
        }
        service, err := youtube.New(client)
        if err != nil {
                log.Printf("error creating youtube client: %v", err)
                return
        }
        call := service.Search.List("id,snippet").Q(msg).MaxResults(1)
        resp, err := call.Do()
        if err != nil {
                log.Printf("error performing youtube search: %v", err)
                return
        }
        var video = ""
        for _,item := range resp.Items {
                switch item.Id.Kind {
                case "youtube#video":
                        video = fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId)
                        break
                }
        }

        actions.Say(video)
}
