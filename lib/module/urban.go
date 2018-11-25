package module

import (
	"fmt"
	"github.com/davidscholberg/go-urbandict"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"os"
	"strings"
)

func HelpUrban() []string {
	s := []string{
		"urban [search phrase] - search urban dictionary for given phrase" +
			" (or get random word if none given)",
	}
	return s
}

func HelpUrbanWotd() []string {
	s := []string{
		"urban_wotd - get the urban dictionary word of the day",
	}
	return s
}

func HelpUrbanTrending() []string {
	s := []string{
		"urban_trending - get the current urban dictionary trending list",
	}
	return s
}

func Urban(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	showSearchResult(in, actions)
}

func UrbanWotd(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	showDefinition(in, actions, true)
}

func UrbanTrending(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	showTrending(in, actions)
}

func showSearchResult(in *message.InboundMsg, actions *Actions) {
	var def *urbandict.Definition
	var err error
	var isRandom bool
	nick := in.Event.Nick
	if len(in.MsgArgs) == 1 {
		def, err = urbandict.Random()
		isRandom = true
	} else {
		def, err = urbandict.Define(strings.Join(in.MsgArgs[1:], " "))
		isRandom = false
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		actions.Say("error: couldn't get search result")
		return
	}

	if isRandom {
		actions.Say(fmt.Sprintf("%s: random word: \"%s\" - %s", nick, def.Word, def.Permalink))
	} else {
		actions.Say(fmt.Sprintf("%s: top result for \"%s\" - %s", nick, def.Word, def.Permalink))
	}
}

func showDefinition(in *message.InboundMsg, actions *Actions, isWotd bool) {
	var def *urbandict.Definition
	var err error
	nick := in.Event.Nick
	if isWotd {
		def, err = urbandict.WordOfTheDay()
	} else if len(in.MsgArgs) == 1 {
		def, err = urbandict.Random()
	} else {
		def, err = urbandict.Define(strings.Join(in.MsgArgs[1:], " "))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		actions.Say("error: couldn't get definition")
		return
	}

	if isWotd {
		actions.Say(fmt.Sprintf("%s: Word of the day: \"%s\"", nick, def.Word))
	} else {
		actions.Say(fmt.Sprintf("%s: Top definition for \"%s\"", nick, def.Word))
	}
	for _, line := range strings.Split(def.Definition, "\r\n") {
		actions.Say(fmt.Sprintf("%s: %s", nick, line))
	}
	actions.Say(fmt.Sprintf("%s: Example:", nick))
	for _, line := range strings.Split(def.Example, "\r\n") {
		actions.Say(fmt.Sprintf("%s: %s", nick, line))
	}
	actions.Say(fmt.Sprintf("%s: permalink: %s", nick, def.Permalink))
}

func showTrending(in *message.InboundMsg, actions *Actions) {
	nick := in.Event.Nick

	trendingWords, err := urbandict.Trending()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		actions.Say("error: couldn't get trending list")
		return
	}

	actions.Say(fmt.Sprintf("%s: Top %d trending words:",
		nick,
		len(trendingWords)))

	for i, word := range trendingWords {
		actions.Say(fmt.Sprintf("%s: %d. %s", nick, i+1, word))
	}

	return
}
