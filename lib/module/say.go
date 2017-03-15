package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"strings"
)

var sayOwner string
var sayDenyMessage string

func ConfigSay(cfg *configure.Config) {
	sayOwner = cfg.Admin.Owner
	sayDenyMessage = cfg.Admin.DenyMessage
}

func HelpSay() []string {
	s := "say <dest> <message> - send message to dest (requires owner privilege)"
	return []string{s}
}

func Say(in *message.InboundMsg, actions *Actions) {
	if in.Event.Nick != sayOwner {
		actions.Say(fmt.Sprintf("%s: %s", in.Event.Nick, sayDenyMessage))
		return
	}

	if len(in.MsgArgs) < 3 {
		msg := "not enough args"
		actions.Say(fmt.Sprintf("%s: %s", in.Event.Nick, msg))
		return
	}

	dest := in.MsgArgs[1]
	msg := strings.TrimSpace(strings.Join(in.MsgArgs[2:], " "))

	actions.SayTo(dest, msg)
}
