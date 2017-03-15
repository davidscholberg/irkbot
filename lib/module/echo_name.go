package module

import (
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"strings"
)

var nick string

func ConfigEchoName(cfg *configure.Config) {
	nick = cfg.User.Nick
}

func EchoName(in *message.InboundMsg, actions *Actions) bool {
	if !strings.HasPrefix(in.Msg, fmt.Sprintf("%s!", nick)) {
		return false
	}
	actions.Say(fmt.Sprintf("%s!", in.Event.Nick))
	return true
}
