package modpm

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib"
	"strings"
)

var helpMsgs []string

// RegisterHelp allows modules to define help strings to be displayed on command.
func RegisterHelp(s []string) {
	helpMsgs = append(helpMsgs, s...)
}

// Help displays help for all bot commands.
func Help(p *lib.Privmsg) bool {
	if !strings.HasPrefix(p.Msg, "..help") {
		return false
	}

	nick := p.Event.Nick

	lib.Say(p, fmt.Sprintf("%s: List of commands:", nick))

	for _, s := range helpMsgs {
		lib.Say(p, fmt.Sprintf("%s: %s", nick, s))
	}

	return true
}
