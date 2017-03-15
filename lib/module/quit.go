package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
)

var owner string
var denyMessage string

func ConfigQuit(cfg *configure.Config) {
	owner = cfg.Admin.Owner
	denyMessage = cfg.Admin.DenyMessage
}

func HelpQuit() []string {
	s := "quit - ragequit IRC (requires owner privilege)"
	return []string{s}
}

func Quit(in *message.InboundMsg, actions *Actions) {
	if in.Event.Nick != owner {
		actions.Say(fmt.Sprintf("%s: %s", in.Event.Nick, denyMessage))
		return
	}
	actions.Quit()
}
