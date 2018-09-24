package module

import (
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
)

func HelpShrug() []string {
	s := "shrug - ¯\\_(ツ)_/¯"
	return []string{s}
}

func Shrug(cfg *configure.Config, in *message.InboundMsg, actions *Actions) {
	actions.Say("¯\\_(ツ)_/¯")
}
