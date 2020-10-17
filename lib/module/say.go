package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"strings"
)

func helpSay() []string {
	s := "say <dest> <message> - send message to dest (requires owner privilege)"
	return []string{s}
}

func say(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	if in.Event.Nick != cfg.Admin.Owner {
		actions.say(fmt.Sprintf("%s: %s", in.Event.Nick, cfg.Admin.DenyMessage))
		return
	}

	if len(in.MsgArgs) < 3 {
		msg := "not enough args"
		actions.say(fmt.Sprintf("%s: %s", in.Event.Nick, msg))
		return
	}

	dest := in.MsgArgs[1]
	msg := strings.TrimSpace(strings.Join(in.MsgArgs[2:], " "))

	actions.sayTo(dest, msg)
}
