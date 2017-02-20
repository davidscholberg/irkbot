package module

import (
	"fmt"
	"github.com/davidscholberg/go-urbandict"
	"github.com/davidscholberg/irkbot/lib/message"
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

func Urban(p *message.Privmsg) {
	showSearchResult(p)
}

func UrbanWotd(p *message.Privmsg) {
	showDefinition(p, true)
}

func UrbanTrending(p *message.Privmsg) {
	showTrending(p)
}

func showSearchResult(p *message.Privmsg) {
	var def *urbandict.Definition
	var err error
	var isRandom bool
	nick := p.Event.Nick
	if len(p.MsgArgs) == 1 {
		def, err = urbandict.Random()
		isRandom = true
	} else {
		def, err = urbandict.Define(strings.Join(p.MsgArgs[1:], " "))
		isRandom = false
	}
	if err != nil {
		message.Say(p, fmt.Sprintf("%s: %s", nick, err.Error()))
		return
	}

	if isRandom {
		message.Say(p, fmt.Sprintf("%s: random word: \"%s\" - %s", nick, def.Word, def.Permalink))
	} else {
		message.Say(p, fmt.Sprintf("%s: top result for \"%s\" - %s", nick, def.Word, def.Permalink))
	}
}

func showDefinition(p *message.Privmsg, isWotd bool) {
	var def *urbandict.Definition
	var err error
	nick := p.Event.Nick
	if isWotd {
		def, err = urbandict.WordOfTheDay()
	} else if len(p.MsgArgs) == 1 {
		def, err = urbandict.Random()
	} else {
		def, err = urbandict.Define(strings.Join(p.MsgArgs[1:], " "))
	}
	if err != nil {
		message.Say(p, fmt.Sprintf("%s: %s", nick, err.Error()))
		return
	}

	// TODO: implement max message length handling

	if isWotd {
		message.Say(p, fmt.Sprintf("%s: Word of the day: \"%s\"", nick, def.Word))
	} else {
		message.Say(p, fmt.Sprintf("%s: Top definition for \"%s\"", nick, def.Word))
	}
	for _, line := range strings.Split(def.Definition, "\r\n") {
		message.Say(p, fmt.Sprintf("%s: %s", nick, line))
	}
	message.Say(p, fmt.Sprintf("%s: Example:", nick))
	for _, line := range strings.Split(def.Example, "\r\n") {
		message.Say(p, fmt.Sprintf("%s: %s", nick, line))
	}
	message.Say(p, fmt.Sprintf("%s: permalink: %s", nick, def.Permalink))
}

func showTrending(p *message.Privmsg) {
	nick := p.Event.Nick

	trendingWords, err := urbandict.Trending()
	if err != nil {
		message.Say(p, fmt.Sprintf("%s: %s", nick, err.Error()))
		return
	}

	message.Say(p, fmt.Sprintf("%s: Top %d trending words:",
		nick,
		len(trendingWords)))

	for i, word := range trendingWords {
		message.Say(p, fmt.Sprintf("%s: %d. %s", nick, i+1, word))
	}

	return
}
