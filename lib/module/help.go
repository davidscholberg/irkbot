package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"sort"
	"strings"
)

var helpMsgs []string

// RegisterHelp allows modules to define help strings to be displayed on command.
func registerHelp(s []string) {
	helpMsgs = append(helpMsgs, s...)
}

func sortHelp() {
	sort.Strings(helpMsgs)
	return
}

// Help displays help for all bot commands.
func help(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	nick := in.Event.Nick

	if strings.HasPrefix(in.Src, "#") {
		actions.say(
			fmt.Sprintf(
				"%s: check your PMs, fam",
				nick,
			),
		)
	}

	actions.sayTo(nick, "Hello! I am an Irkbot instance - "+
		"https://github.com/davidscholberg/irkbot")
	actions.sayTo(nick, "Here's my list of commands:")

	for _, s := range helpMsgs {
		actions.sayTo(nick, fmt.Sprintf("%s%s", cfg.Channel.CmdPrefix, s))
	}
}

func parseHelp(cfg *configure.Config, in *message.InboundMsg, actions *actions) bool {
	if strings.TrimSpace(in.Msg) != fmt.Sprintf("%s: help", cfg.User.Nick) {
		return false
	}

	help(cfg, in, actions)

	return false
}
