package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
)

var cmdPrefix string
var helpMsgs []string

func ConfigHelp(cfg *configure.Config) {
	cmdPrefix = cfg.Channel.CmdPrefix
}

// RegisterHelp allows modules to define help strings to be displayed on command.
func RegisterHelp(s []string) {
	helpMsgs = append(helpMsgs, s...)
}

// Help displays help for all bot commands.
func Help(in *message.InboundMsg, actions *Actions) {
	nick := in.Event.Nick

	actions.Say(fmt.Sprintf("%s: Hello! I am an Irkbot instance - "+
		"https://github.com/davidscholberg/irkbot", nick))
	actions.Say(fmt.Sprintf("%s: Here's my list of commands:", nick))

	for _, s := range helpMsgs {
		actions.Say(fmt.Sprintf("%s: %s%s", nick, cmdPrefix, s))
	}
}
