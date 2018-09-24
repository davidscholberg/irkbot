package module

import (
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
)

func HelpLenny() []string {
	s := "lenny - ( ͡° ͜ʖ ͡°)"
	return []string{s}
}

func Lenny(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	actions.Say("( ͡° ͜ʖ ͡°)")
}
