package module

import (
	"fmt"
	"github.com/davidscholberg/go-urbandict"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"os"
	"strings"
)

func helpUrban() []string {
	s := []string{
		"urban [search phrase] - search urban dictionary for given phrase" +
			" (or get random word if none given)",
	}
	return s
}

func helpUrbanWotd() []string {
	s := []string{
		"urban_wotd - get the urban dictionary word of the day",
	}
	return s
}

func helpUrbanTrending() []string {
	s := []string{
		"urban_trending - get the current urban dictionary trending list",
	}
	return s
}

func urban(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	showSearchResult(in, actions)
}

func urbanWotd(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	showDefinition(in, actions, true)
}

func urbanTrending(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	showTrending(in, actions)
}

func showSearchResult(in *message.InboundMsg, actions *actions) {
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
		actions.say("error: couldn't get search result")
		return
	}

	if isRandom {
		actions.say(fmt.Sprintf("%s: random word: \"%s\" - %s", nick, def.Word, def.Permalink))
	} else {
		actions.say(fmt.Sprintf("%s: top result for \"%s\" - %s", nick, def.Word, def.Permalink))
	}
}

func showDefinition(in *message.InboundMsg, actions *actions, isWotd bool) {
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
		actions.say("error: couldn't get definition")
		return
	}

	if isWotd {
		actions.say(fmt.Sprintf("%s: Word of the day: \"%s\"", nick, def.Word))
	} else {
		actions.say(fmt.Sprintf("%s: Top definition for \"%s\"", nick, def.Word))
	}
	for _, line := range strings.Split(def.Definition, "\r\n") {
		actions.say(fmt.Sprintf("%s: %s", nick, line))
	}
	actions.say(fmt.Sprintf("%s: Example:", nick))
	for _, line := range strings.Split(def.Example, "\r\n") {
		actions.say(fmt.Sprintf("%s: %s", nick, line))
	}
	actions.say(fmt.Sprintf("%s: permalink: %s", nick, def.Permalink))
}

func showTrending(in *message.InboundMsg, actions *actions) {
	nick := in.Event.Nick

	trendingWords, err := urbandict.Trending()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		actions.say("error: couldn't get trending list")
		return
	}

	actions.say(fmt.Sprintf("%s: Top %d trending words:",
		nick,
		len(trendingWords)))

	for i, word := range trendingWords {
		actions.say(fmt.Sprintf("%s: %d. %s", nick, i+1, word))
	}

	return
}
