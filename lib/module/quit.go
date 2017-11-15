package module

import (
	"fmt"
	"github.com/jholtom/irkbot/lib/configure"
	"github.com/jholtom/irkbot/lib/message"
)

func HelpQuit() []string {
	s := "quit - ragequit IRC (requires owner privilege)"
	return []string{s}
}

func Quit(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	if in.Event.Nick != cfg.Admin.Owner {
		actions.Say(fmt.Sprintf("%s: %s", in.Event.Nick, cfg.Admin.DenyMessage))
		return
	}
	actions.Quit()
}
