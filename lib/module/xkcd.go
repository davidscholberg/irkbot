package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/nishanths/go-xkcd"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const apiUrlFmtDefine = "https://relevantxkcd.appspot.com/process?%s"

func Helpxkcd() []string {
	s := "xkcd <search> - find a xkcd comic relevant to the search term"
	return []string{s}
}
func get(apiURL string) (string, error) {
	response, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	return bodyString, err
}
func fmtString(bodyString string) string {
	bodyStrings := strings.Split(bodyString, "\n")
	spacedStrings := strings.Fields(bodyStrings[2])
	return spacedStrings[0]
}
func getXKCD(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	comic := "enter a search term, asshat"
	//If no search term, gently remind the user to input one
	if len(in.MsgArgs[1:]) == 0 {
		actions.Say(comic)
		return
	}
	query := url.Values{}
	query.Add("action", "xkcd")
	search := strings.Join(in.MsgArgs[1:], " ")
	query.Add("query", search)
	apiUrl := fmt.Sprintf(apiUrlFmtDefine, query.Encode())
	comicString, comicErr := get(apiUrl)
	if comicErr != nil {
		actions.Say("something borked, try again")
		return
	}
	comicNum := fmtString(comicString)
	client := xkcd.NewClient()
	i, strconvErr := strconv.Atoi(comicNum)
	if strconvErr != nil {
		actions.Say("something borked, try again")
		return
	}
	comicGet, err := client.Get(i)
	if err != nil {
		actions.Say("something borked, try again")
		return
	}
	comic = fmt.Sprintf("https://xkcd.com/%s/ - %s ", comicNum, comicGet.Title)
	actions.Say(comic)
}
