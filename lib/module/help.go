package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"strings"
)

var helpMsgs []string

// RegisterHelp allows modules to define help strings to be displayed on command.
func RegisterHelp(s []string) {
	helpMsgs = append(helpMsgs, s...)
}

// Help displays help for all bot commands.
func Help(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	nick := in.Event.Nick

	if strings.HasPrefix(in.Src, "#") {
		actions.Say(
			fmt.Sprintf(
				"%s: check your PMs, fam",
				nick,
			),
		)
	}

	actions.SayTo(nick, "Hello! I am an Irkbot instance - "+
		"https://github.com/davidscholberg/irkbot")
	actions.SayTo(nick, "Here's my list of commands:")

	for _, s := range helpMsgs {
		actions.SayTo(nick, fmt.Sprintf("%s%s", cfg.Channel.CmdPrefix, s))
	}
}
