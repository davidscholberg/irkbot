package module

import (
	"fmt"
	"github.com/dvdmuckle/irkbot/lib/configure"
	"github.com/dvdmuckle/irkbot/lib/message"
	"strings"
)

func HelpSay() []string {
	s := "say <dest> <message> - send message to dest (requires owner privilege)"
	return []string{s}
}

func Say(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	if in.Event.Nick != cfg.Admin.Owner {
		actions.Say(fmt.Sprintf("%s: %s", in.Event.Nick, cfg.Admin.DenyMessage))
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
