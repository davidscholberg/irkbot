package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/nishanths/go-xkcd/v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const apiUrlFmtDefine = "https://relevantxkcd.appspot.com/process?%s"

func helpxkcd() []string {
	s := "xkcd <search> - find a xkcd comic relevant to the search term"
	return []string{s}
}

//Perform the actual GET and return the resulting body as a string.
//This function closes the response body.
func getComicString(response *http.Response, responseSizeLimit int64) (string, error) {
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(io.LimitReader(response.Body, responseSizeLimit))
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	return bodyString, nil
}

//Parse the body string for the comic number we want
func parseString(bodyString string) (string, error) {
	bodyStrings := strings.Split(bodyString, "\n")
	if len(bodyStrings) < 3 {
		return "", fmt.Errorf("error in parsing string: splitting body by line failed")
	}
	spacedStrings := strings.Fields(bodyStrings[2])
	if len(spacedStrings) < 1 {
		return "", fmt.Errorf("error in parsing string: accessing substring of bodyStrings failed")
	}
	return spacedStrings[0], nil
}

//Method called on xkcd command, named funky so as not to collide with xkcd-go
func getXKCD(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	comicMsg := "enter a search term plz"
	//If no search term, gently remind the user to input one
	if len(in.MsgArgs[1:]) == 0 {
		actions.say(comicMsg)
		return
	}
	query := url.Values{}
	query.Add("action", "xkcd")
	search := strings.Join(in.MsgArgs[1:], " ")
	query.Add("query", search)
	apiUrl := fmt.Sprintf(apiUrlFmtDefine, query.Encode())
	response, err := actions.httpGet(apiUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		actions.say("something borked, try again")
		return
	}
	comicString, comicErr := getComicString(response, cfg.Http.ResponseSizeLimit)
	if comicErr != nil {
		fmt.Fprintln(os.Stderr, comicErr)
		actions.say("something borked, try again")
		return
	}
	comicNum, parseErr := parseString(comicString)
	if parseErr != nil {
		fmt.Fprintln(os.Stderr, parseErr)
		actions.say("something borked, try again")
		return
	}
	client := &xkcd.Client{
		HTTPClient: &http.Client{Timeout: time.Duration(cfg.Http.Timeout) * time.Second},
		Config:     xkcd.Config{UseHTTPS: true},
	}
	i, strconvErr := strconv.Atoi(comicNum)
	if strconvErr != nil {
		fmt.Fprintln(os.Stderr, strconvErr)
		actions.say("something borked, try again")
		return
	}
	comicGet, err := client.Get(i)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		actions.say("something borked, try again")
		return
	}
	comicMsg = fmt.Sprintf("%s - https://xkcd.com/%s/", comicGet.Title, comicNum)
	actions.say(comicMsg)
}
