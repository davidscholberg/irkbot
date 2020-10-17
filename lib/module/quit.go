package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
)

func helpQuit() []string {
	s := "quit - ragequit IRC (requires owner privilege)"
	return []string{s}
}

func quit(cfg *configure.Config, in *message.InboundMsg, actions *actions) {
	if in.Event.Nick != cfg.Admin.Owner {
		actions.say(fmt.Sprintf("%s: %s", in.Event.Nick, cfg.Admin.DenyMessage))
		return
	}
	actions.quit()
}
