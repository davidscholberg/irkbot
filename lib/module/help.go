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
func Help(p *message.Privmsg) {
	nick := p.Event.Nick

	message.Say(p, fmt.Sprintf("%s: List of commands:", nick))

	for _, s := range helpMsgs {
		message.Say(p, fmt.Sprintf("%s: %s%s", nick, cmdPrefix, s))
	}
}
