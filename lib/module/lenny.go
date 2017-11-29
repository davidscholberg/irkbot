package module

import (
	"github.com/dvdmuckle/irkbot/lib/configure"
	"github.com/dvdmuckle/irkbot/lib/message"
)

func HelpLenny() []string {
	s := "lenny - you know what this does"
	return []string{s}
}

func Lenny(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	actions.Say("( ͡° ͜ʖ ͡°)")
}
